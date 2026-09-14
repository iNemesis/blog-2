---
title: Manuel du tech lead efficace
date: 2026-09-12
tags:
  - tech lead
---

**/!\ cet article est un premier jet**

J’en ai marre des tech leads mous.

J'ai donc identifié trois piliers qui regroupent les bonnes pratiques de tech leading qui fonctionnent pour moi, et pour les tech leads qui m'entourent.


> “You will learn to have uncomfortable conversations. You will learn to resolve conflicts, dear god you ever learn to resolve conflicts. Actually you’ll learn to yearn for conflicts because straightforward conflict is usually better than all the other options.” – Charity.wtf

**Premièrement, vous devez être un meneur**. L’expérience technique est nécessaire mais secondaire par rapport au leading.

Pour mener, vous devez donner des directives claires, des défis, vous devez faire des feedbacks, montrer l'exemple, et surtout assumer la responsabilités des décisions de l’équipe.


> “Great teams are unafraid to air their dirty laundry. They admit their mistakes, their weaknesses, and their concerns without fear of reprisal” – Patrick Lencioni

**Deuxièmement, vous devez servir l’équipe**. Vous devez abattre les obstacles sur sa route et mettre de l'huile dans les rouages.

Vous devez protéger le flow de l'équipe en allant en réunions pour eux, ce qui vous permettra aussi de cadrer(1 définition) les sujets et de régler les dépendances en amont. Ensuite, vous devez créer un environnement où chaque membre peut s'exprimer, faire des erreurs et poser des questions "naïves". Vous devez détecter les signaux faibles de vos équipiers, et les aider à régler leurs problèmes avant qu'ils ne pensent à partir. Vous devez aussi les challenger pour les faire progresser, en leur donnant la cane à pêche plutôt que le poisson, donc en les laissant échouer.


> “C’est la compréhension du métier par les développeurs qui part en production.” – Alberto Brandolini
 
**Troisièmement, vous devez concevoir un produit qui correspond aux besoins des métiers et des utilisateurs.** Le tout en restant maintenable.

Pour cela, vous devez créer une synergie avec votre PO pour comprendre le métier et faire comprendre les contraintes techniques. Vous devez supprimer les frontières qui se trouvent entre les experts métiers et votre équipe, pour impliquer les développeurs et les faire passer d'exécutants à créateurs. Pour finir, vous devez fédérer l’équipe autour d’un objectif métier commun.


Être tech lead n’est pas simple, et ce chemin n’est pas permis à tous les développeurs, aussi experts soient-ils. C'est pas le tech de l'équipe qui doit passer tech lead, mais celui qui a l'étoffe d'un leader.



## I Mener l'équipe

> "There are two types of people: those who try to win and those who try to win arguments. They are never the same." - Nassim Nicholas Taleb

### Donner des enjeux
Challenger l'équipe en augmentant les objectifs, d'être ambitieux. Par exemple, on a récemment démarré un nouveau projet et le client s'attendait à ce qu'on fasse 3 fonctionnalités sur 2 mois. En équipe, j'ai poussé pour qu'on s'engage sur une fonctionnalité en plus, pour qu'on soit dans un flow soutenu. L'équipe a accepté et on a réussi à faire encore plus que ce sur quoi on s'était engagé.

L'objectif est d'être dans le sweet spot du flux. Pour cela j'utilise la loi de Parkinson "tout travail augmente jusqu’à occuper entièrement le temps qui lui est affecté". Le but est de choisir le juste milieu entre burnout et boreout (l'inverse de burnout, quand l'équipe s'ennuie). C'est à dire faire plus de choses en moins de temps.

// ajouter schéma

### Déléguer
En tant que TL, vous devez donner des directives, vous devez être direct. Vous ne devez pas avoir froid aux yeux de demander des choses, vous êtes là pour ça. Le meilleur moyen de donner des directives est de poser des questions : "Tu pourrais regarder ce ticket ?", "tu aurais le temps de demander à l'équipe X comment gérer l'authentification ?"...

En 2025, je suis repassé développeur quelques mois. Mon tech lead était dans son premier poste, il n'osait pas donner des directives à l'équipe. Nous étions chacun des électrons libres qui nous concentrions sur le backlog. Si des problèmes survenaient en dehors du backlog (dépendance avec une autre équipe, maintient du produit, problèmes de prod), nous étions exclus. Il n'y avait pas de partage.

A ce moment là, j'aurais aimé que mon TL me dise "il y a un pb en prod, vas-y" ou "gère cette dépendance avec l'équipe" ou "prend ce sujet". Sans ordre clair, l'équipe prend le sujet le plus important pour elle, qui n'est pas toujours le plus utile pour le produit.

### Faire des feedbacks
> On souffre plus du manque de feedback que des feedbacks rudes. 

Vous ne devez jamais faire de feedbacks négatifs en compagnie d'autres personnes. Les feedbacks se font toujours en privé.

J'utilise le framework "BEEF", Behaviour, Examples, Effects et Future. "Je dois te parler de tes retards aux daily. Tu étais en retard de 10 minutes ce matin et tu n'es pas venu vendredi dernier. Tout l'équipe t'a attendu et le client s'est posé des questions. J'aimerais que tu nous préviennes en amont lorsque tu as un retard, et que tu ne sois pas en retard par défaut.".

Ces feedbacks doivent vous appartenir "J'ai l'impression que...", "J'ai trouvé que...". Personne ne peut vous empêcher de ressentir. Mais dire "tu étais énérvé", l'interlocuteur peut dire qu'il ne l'était pas et ça peut le renfermer.

Vous n'avez pas besoin de suivre un framework pour faire des feedbacks, ça ne doit pas vous empêcher d'en faire. Mais si vous devez en faire, évitez à tout prix les feedbacks sandwich "tu es très sympa, mais tu arrives tard le matin, et tu travailles super bien.". Vous diluez votre feedback tout en vous faisant passer pour un manager nul.

Les feedbacks ne sont pas que négatifs, il faut aussi que vous fassiez des feedbacks positifs à votre équipe. Un ticket réalisé avec succès, des questions pertinentes, des remarques sur l'architecture... Vous devez donner du retour sur ce que l'équipe fait de bien.

Vous devez aussi faire résonner ces feedbacks. Par exemple, je vais parfois voir les managers des développeurs de mon équipe pour leur dire qu'ils mériteraient d'être récompensés.

#### Sans rancunes
Parfois, les feedbacks se passent mal. Vous pouvez dire un mot qui *déclenche* votre interlocuteur. Ca m'est arrivé récemment, j'ai dit à un collègue qu'il silotait son équipe car à chaque fois qu'on avait besoin d'eux, on devait se justifier. On passait parfois plus de temps à nous justifier qu'à résoudre notre problème. Je trouvais que c'était dommage et quand une personne de mon équipe m'a dit la même chose, ça a validé mon intuition et je suis allé lui parler. Il a pété un câble. Il m'a dit qu'il protégeait son équipe, et j'ai eu une parole de trop quand j'ai dit "c'est bon, ils ne sont pas en sucre", ça l'a énérvé. Le débat était terminé. Il contredisait tout ce que je lui disait "vous êtes 12, dans votre équipe" "non" "combien?" "on est 11". Je parlais à un mur. J'ai donc arrêté la conversation.

Parfois, on dit des mots de trop. Vous n'êtes pas un acteur avec un dialogue préparé d'avance. Vous allez peut être heurté votre interlocteur sans même vous en rendre compte. L'essentiel est d'arrêter la conversation et d'attendre que l'autre personne redescende en émotion. 

De votre côté, vous ne devez pas lui en tenir rigueur. On a tous des émotions, et on a parfois du mal à les contenir. J'ai mis mon égo de côté et je suis allé le revoir le lendemain matin pour lui dire bonjour. Sans rancune.

Ca ne veut pas dire que vous devez être ami avec tout le monde. Il y a des personnes avec qui je ne veux pas travailler, et je fais tout pour trouver d'autres interlocuteurs. Mais parfois, on doit tout de même travailler avec. Vous devez mettre votre égo et vos préférences de côté. L'important c'est d'aller de l'avant. Le produit doit avancer.

On parle souvent d'egoless programming, mais je pense qu'il ne faut pas limiter ce principe au développement. Vous devez faire de l'egoless working. Travailler sans égo.

#### Celui qui s'énerve a perdu
Dans une conversation, celui qui s'énerve a perdu. Et vous allez vous énerver. Un jour, il fera trop chaud dans la salle ou vous aurez passé une mauvaise journée. Vous allez avoir une conversation et vous n'arriverez pas à vous contenir. 

L'important c'est d'avancer. Si vous vous énervez, vous devenez un mur qui n'avance pas. Vous pouvez même faire reculer vos relations, et donc le projet. Donc dans ce cas, partez. Vous reviendrez le lendemain avec un meilleur mood et du recul sur le débat.

### Montrer l'exemple
> Rien n'est plus contagieux que l'exemple. - François de La Rochefoucauld

Pour que votre équipe grandisse et vous écoute, vous devez vivre parmi eux et incarner les valeurs que vous défendez. En incarnant les valeurs que vous souhaitez voir dans l'équipe, vos équipiers vont prendre exemple sur vous.

Pour cela, vous devez être sur le terrain, vous devez coder. On parle souvent de passer au minimum 30% de votre temps dans le code. Cela vous permet de comprendre les douleurs de vos équipiers, de pouvoir proposer de vraies améliorations de design, et surtout d'être plus pertinent face aux acteurs métier.

Vous devez donc être assez dans le code pour le comprendre, et vous devez aller assez en réunion pour comprendre le métier et le challenger.

## II Servir l'équipe

### Faire des O3
Ces O3 me permettent de détecter les signaux faibles. Par exemple, si l'un de mes développeur m'a l'air moins motivé que d'habitude, je peux lui en parler lors de ces points.

Je veux savoir si mes développeurs vont bien et résoudre leurs problèmes avant que ceux-ci soient trop gros pour être résolu, et qu'ils quittent l'équipe. Une équipe qui dure est une équipe productive.

En réalité, ces points parlent surtout des aspirations de l'équipe. Par exemple, c'est à ce moment là qu'un de mes développeur me dit qu'il veut faire de l'OPS ou qu'il s'intéresse à Keycloak. S'ils se projettent pour devenir Tech lead, je peux alors les guider vers ce chemin et leur donner du leste sur mon rôle.

Ces moments d'échange sont très appréciés par l'équipe.

#### En marchant
J'ai repris l'une des habitudes de Steve Jobs pour ces O3, je les fait en marchant. Et j'y ai découvert trois avantages. Le premier est que le côte cringe d'un O3 disparait, le second est que ça devient une activité plaisante, pour prendre un bain de soleil, et la troisième est que ça libère la parole. En effet, quand le corps est occupé, l'esprit est concentré et permet des conversations plus profondes.

Dans une de mes équipes, ça faisait des semaines que je faisais des O3 avec notre PO et je sentais que quelque chose n'allait pas de son côté. Mais j'étais face à un mur. Lors de nos O3, tout allait bien selon elle. Je n'arrivais pas à valider mes intuitions. Etant en train de lire la biographie de Steve Jobs par Walter Isaacson, j'ai lu cette pratique de les faire en marchant. Ca a tout changé. Elle s'est ouverte et m'a confié pourquoi ça n'allait pas. Ce qui m'a permis d'intervenir.

Je pratique ces points une fois par mois par équipier. Je poses des questions ouvertes comme "Quelle compétence voudrais-tu développer ?", "Comment te sens-tu dans l'équipe ?", "Qu'as-tu pensé de \[évènement] ?", "Qu'est-ce qu'on pourrait supprimer pour livrer plus vite ?". Et surtout, je marche avec eux, en ce moment au parc de Bercy.


### Donner du leste
Vous devez donner du leste à vos équipes. Vous devez les laisser faire des erreurs, même si vous savez que c'est une erreur. L'un de mes profs de fac disait "\[Dans l'apprentissage] Vaut mieux guérir que prévenir". Vous pouvez expliquer des concepts pendant deux semaines, mais vous ne serez jamais autant efficace qu'en les laissant pratiquer et échouer.

En 2024, alors que j'arrivais dans une nouvelle équipe et découvrais encore le produit, l'un des développeur souhaitait factoriser trois formulaires qui se ressemblaient beaucoup en un seul. Par défaut, je suis assez réticent à ce genre de factorisation. Ces formulaires avaient presque tous les champs en commun mais représentaient tous trois des concepts métiers différents. De plus, nous venions de les développer, et ils n'avaient pas encore passés leur baptême du feu : les retours clients.

Mais je venais d'arriver dans l'équipe, et même avec mes meilleurs arguments et exemples, je voyais que ce développeur n'était pas convaincu. J'aurais pu m'opposer à cette factorisation mais je pense que ce reflexe serait revenu plus tard ou dans un autre composant.

Alors j'ai préféré lui accorder cette factorisation, et la faire avec lui. Après tout, ce composant n'était pas sensible, et n'était pas sur un flux critique.

Un mois et demi plus tard, le composant avait évolué. On avait ajouté des contraintes de validation propres à chaque concept métier. Et ce composant qui était concis et simple à l'époque est devenu un monstre de complexité. Il contenait des *if* partout pour savoir dans quel concept métier il était utilisé. Le concept SOLID de SRP était brisé.

Donc on a redécoupé ce composant en trois, comme au début. Ca a pris deux jours mais le code est redevenu lisible et évolutif. Et surtout, le développeur qui voulait tout factoriser a appris que la factorisation n'est pas toujours la bonne solution, et que dupliquer c'est découpler.

Par contre, vous ne pouvez pas laisser du leste partout. Vous devez identifier les parties du code non sensibles sur lesquels votre équipe peut apprendre par l'erreur. Restez le garde fou tout en restant ouvert.

### Créer un environnement sain
Pour récupérer les meilleures idées et remarques de chacun, dans le but de créer le meilleur produit, il faut que tout le monde puisse s'exprimer comme il le souhaite. Pour cela, il faut un environnement sain, où chacun se sent en sécurité.

Je montre l'exemple en posant des questions bêtes, en acceptant quand j'ai tort et que mes équipiers ont raison. Ca n'est plus un sujet d'avoir tort. 

Je m'assure aussi que tout le monde se soit exprimé en distribuant la parole. Si j'ai remarqué qu'un de mes équipier allait parler mais que quelqu'un d'autre a pris le flux de conversation, je lui redistribue à la fin. 

Je m'assure aussi que les sujets soient traités un par un. Le but est d'explorer chaque idée en profondeur. 

Mais le plus important pour créer un environnement sain est d'accepter les erreurs. C'est grâce aux erreurs que l'équipe apprend. Vous devez être présent pour aider à les réparer, pour faire en sorte qu'elles ne se produisent pas, et surtout pour en prendre la responsabilité.

### Aller en réunion pour protéger le flow
Le tech lead va à plus de réunions que ses développeurs. C'est d'ailleurs l'une des principales douleurs lorsqu'un développeur passe tech lead.

Mais le tech lead va en réunion pour la bonne cause, il y va pour que ses développeur puissent coder.

Parce qu'en allant dans ces réunions, il représente l'équipe, et il règle les dépendances avec les autres équipes, qui sont des obstacles au flow.

#### Protéger son propre agenda
Par principe, je ne dis jamais NON à l'équipe ni aux experts métier. Je peux éventuellement challenger si je doute de l'utilité d'un point mais s'ils en ont besoin, j'y vais.

Par contre, pour les managers du projet, par exemple les product manager ou les architectes, j'ai plus de mal à dire oui. Assez-souvent, ces points ne me servent pas, ils leur servent à eux, pour savoir ce que l'équipe fait. Ils agissent comme des passes plats pour la direction. Je leur propose donc de leur faire un compte rendu écrit.

J'ajoute aussi que je me permet de ne pas aller à ces points car je remonte des alertes et je n'ai aucun souci pour dire que je ne sais pas ou que je vais avoir besoin d'eux. Certains de leurs points sont là au cas où les équipes ne remontent pas ces informations.

#### L'agenda des managers n'est pas celui des développeur
Selon Paul Graham, il existe deux types d'agenda : celui des managers qui se découpe finement et celui des créateurs qui se découpe en demi journée. Les développeurs sont des créateurs, si vous leur mettez un point à 15 heure, vous gâchez leur après-midi. Pareil, si vous leur mettez un point à 16h30, leur état de concentration sera moindre car il occupera une partie de leur cerveau.

Vous devez faire attention aux managers qui ne se rendent pas compte qu'une heure en plein milieu de journée d'un développeur ne coûte pas juste une heure mais la demi-
journée.

Je vous conseille donc de privilégier les points en bordure de journée, de rassembler les points les jours en présentiel, de faire des points de 10 ou 20 minutes par défaut, et surtout, d'aller en point pour eux.
## III Construire un produit efficace

### Challenger le métier

> "Techniquement, tout est faisable"

C'est très important pour un tech lead d'être le garde fou des fonctionnalités inutiles. Au plus vous allez ajouter de fonctionnalités, même très simples, dans votre application, et au plus celle-ci va devenir complexe. Le nombre de fonctionnalité d'une application est corrélé à la complexité de celle-ci. Vous aurez plus de fonctionnalité à maintenir, plus de fonctionnalité à montrer aux nouveaux arrivants, plus de fonctionnalité à comprendre avant d'en ajouter d'autres. La charge cognitive de votre équipe est corrélée au nombre de fonctionnalité et à leur complexité. Less is more.

> « Celui qui ne peut pas décrire le problème ne trouvera jamais la solution. » -- Confucius

Poser des questions permet de comprendre pourquoi les experts métiers, architectes, autres équipes... veulent ce qu'ils demandent. La question la plus puissante que vous puissiez poser est "Pourquoi ?", au risque de passer pour quelqu'un qui ne comprend pas. Vous devez poser des questions qui vous paraissent bêtes.

Par exemple, les experts métier de l'un de mes projet voulaient une recherche avec une dizaine de champs différents. J'ai donc demandé "Pourquoi avez-vous besoin de ces 10 champs ?", à quoi ils ont répondu "Pour faire des statistiques". Ils étaient habitués à ce mode de fonctionnement car leur précédente application ne fournissait pas de statistiques. J'ai donc challengé ce besoin en leur demandant de ne garder que deux champs. Très vite, ils ont négocié pour 3 champs et ne sont plus jamais revenus sur les 7 autres champs.

### Objectif métier commun
La seconde qui dure dans le temps mais qui est complexe à trouver est de trouver un enjeux métier sur lequel l'équipe a un impact. Par exemple, dans l'une de mes équipe, c'était une dépendance avec un progiciel à remplacer. Sans cet objectif, l'entreprise qui m'engageait perdait plusieurs pourcents de son chiffre d'affaire. On ne pouvait pas échouer.

Mais lorsque cet objectif s'est terminé, l'équipe s'est relachée, à la limite du boreout. Comme le "choc de la retraite", qui veut qu'on tombe malade quand notre rythme de travail s'arrête. Les boreout sont autant néfaste pour l'équipe que les burnout. Cette période calme, survenue juste après une période intense, a amené deux développeurs a quitter l'équipe.

Alors pour remotiver l'équipe, on a créé un dashboard avec des données métier, le prix moyen d'un panier, quel flux d'argent avait généré l'application, le nombre d'utilisateurs... Nos développements avaient un vrai impact sur ces métriques. Par exemple, lorsqu'un bug n'affichait plus certains produits au catalogue, on voyait une perte nette sur la semaine. On avait de vrais enjeux, sur lesquels on avait du pouvoir.


### Transformer vos intuitions en métriques

Pour piloter votre codebase et vous assurer de sa maintenabilité, n'attendez pas un audit externe à 50k, auditez la vous-même. 

N'abusez pas non plus des métriques génériques (coucou SonarQube) qui vous inventent des problèmes (dupliquer c'est découpler). Je vous conseille plutôt de partir des intuitions de votre équipe pour créer des métriques. 

Ne vous prenez pas trop la tête, utilisez des métriques simples (par exemple le nombre de lignes). Grâce à ces métriques, vous pourrez expliquer aux décideurs (coucou les archis) que quelque chose ne va pas et qu'il vous faut du temps pour le corriger. Ces métriques vous pointeront aussi ce qui ne va pas et vous aiguilleront sur le "pourquoi ça ne va pas". 

De plus, en combinant ces métriques, vous pourrez savoir lesquelles valent le coup. Par exemple en couplant votre matrice de dépendance sur Intellij et les fichiers les plus modifiés via une commande git, vous saurez quels sont les fichiers à problèmes vous aurez le + d'impact.

Sur l'un de mes anciens projets, l'équipe passait sa vie à écrire des tests. On ne pouvait pas écrire une ligne sans en casser une douzaine, et le plus difficile n'était pas dans le code mais dans leur écriture. Alors pour valider notre intuition "on teste trop", on a comparé le nombre de lignes de code au nombre de lignes de test. On avait 9 lignes de test par ligne de code. Sur nos autres projets, ce ratio oscille plutôt entre 1 et 2.

Grâce à cette métrique simple, nous avons pu valider notre intuition et aller la défendre au près des architectes, pour changer notre stratégie de test dont ils ne se rendaient pas compte, car ils ne savaient pas qu'on peut (très) mal tester.


### Réduire les frontières avec les experts métier
Les développeur ne veulent qu'une chose : avoir un impact sur le produit. J'ai connu des équipes qui avaient un lot de ticket à développer, et ils s'exécutaient.

Techniquement, tout est possible, ça dépend du prix que vous souhaitez payer pour l'avoir. Le but du tech lead et de l'équipe est de maximiser le rapport valeur utilisateur / temps passé. Au plus il y a de fonctionnalité, au plus l'application est difficile à maintenir et à faire évoluer.

### Duo TL/PO
En tant que tech lead, vous devez former un duo avec votre Product Owner. Vous devez discuter ensemble des fonctionnalité et prendre des décisions tous les deux pour créer un meilleur produit. Bien sur, vous devrez écouter les UX, les experts métier et autres. Mais l'essentiel est de créer un produit en synergie entre technique et métier.
