package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	loadEnv(".env")

	host := os.Getenv("SFTP_HOST")
	user := os.Getenv("SFTP_USER")
	pass := os.Getenv("SFTP_PASSWORD")
	keyPath := os.Getenv("SFTP_KEY_PATH")
	remoteBase := os.Getenv("SFTP_DIR")

	if host == "" || user == "" || (pass == "" && keyPath == "") {
		fmt.Fprintln(os.Stderr, "Erreur : SFTP_HOST, SFTP_USER et (SFTP_PASSWORD ou SFTP_KEY_PATH) doivent être définis dans .env")
		os.Exit(1)
	}

	port := os.Getenv("SFTP_PORT")
	if port == "" {
		port = "22"
	}
	if !strings.Contains(host, ":") {
		host = host + ":" + port
	}

	if remoteBase == "" {
		remoteBase = "/"
	}

	hostKey := os.Getenv("SFTP_HOST_KEY")
	if hostKey == "" {
		fmt.Fprintln(os.Stderr, "Erreur : SFTP_HOST_KEY manquant dans .env — récupère la clé avec : ssh-keyscan -p <port> <host>")
		os.Exit(1)
	}
	pinned, err := parseHostKey(hostKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur SFTP_HOST_KEY : %v\n", err)
		os.Exit(1)
	}

	var authMethods []ssh.AuthMethod
	if keyPath != "" {
		keyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lecture clé SSH : %v\n", err)
			os.Exit(1)
		}
		var signer ssh.Signer
		if passphrase := os.Getenv("SFTP_KEY_PASSPHRASE"); passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(keyBytes)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur parsing clé SSH : %v\n", err)
			os.Exit(1)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	if pass != "" {
		authMethods = append(authMethods, ssh.Password(pass))
	}

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback(pinned),
		Timeout:         15 * time.Second,
	}

	fmt.Printf("Connexion SFTP à %s (utilisateur : %s)...\n", host, user)
	sshConn, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur connexion SSH : %v\n", err)
		os.Exit(1)
	}
	defer sshConn.Close()

	client, err := sftp.NewClient(sshConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur initialisation SFTP : %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("Connecté avec succès.")

	if err := client.MkdirAll(remoteBase); err != nil {
		fmt.Fprintf(os.Stderr, "Attention création dossier distant : %v\n", err)
	}

	var count int
	var totalBytes int64
	startTime := time.Now()

	err = filepath.WalkDir("docs", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel("docs", p)
		if err != nil || rel == "." {
			return nil
		}

		target := path.Join(remoteBase, filepath.ToSlash(rel))

		if d.IsDir() {
			return client.MkdirAll(target)
		}

		size, err := uploadFile(client, p, target)
		if err != nil {
			return fmt.Errorf("échec envoi %s : %w", p, err)
		}

		count++
		totalBytes += size
		fmt.Printf("✓ %s → %s (%s)\n", p, target, formatSize(size))
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors du déploiement : %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDéploiement terminé en %s : %d fichier(s) envoyé(s) (%s)\n",
		time.Since(startTime).Round(time.Millisecond), count, formatSize(totalBytes))
}

func uploadFile(c *sftp.Client, localPath, remotePath string) (int64, error) {
	src, err := os.Open(localPath)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := c.Create(remotePath)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func loadEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

// hostKeyCallback vérifie la clé d'hôte du serveur contre celle épinglée dans .env.
func hostKeyCallback(pinned ssh.PublicKey) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if key.Type() == pinned.Type() && bytes.Equal(key.Marshal(), pinned.Marshal()) {
			return nil
		}
		return fmt.Errorf("clé d'hôte de %s inattendue (%s) — vérifie SFTP_HOST_KEY dans .env",
			hostname, ssh.FingerprintSHA256(key))
	}
}

// parseHostKey accepte une ligne ssh-keyscan ("host type clé") ou "type clé".
func parseHostKey(s string) (ssh.PublicKey, error) {
	fields := strings.Fields(strings.TrimSpace(s))
	if len(fields) > 2 {
		fields = fields[len(fields)-2:]
	}
	if len(fields) < 2 {
		return nil, fmt.Errorf("format attendu : \"type clé\" (sortie de ssh-keyscan)")
	}
	key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.Join(fields, " ")))
	return key, err
}

func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
