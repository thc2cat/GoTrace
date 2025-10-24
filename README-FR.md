
# GoTrace, a Go Network Analyzer , outil de diagnostic pour la latence et la perte de paquets 🌐

Méta-description : "An open-source CLI tool developed in Go to analyze network performance. Measure latency, packet loss, and trace data paths with this network analyzer."

Ce projet est un outil en ligne de commande développé en Go pour analyser la connectivité et les performances d'un réseau. Il combine les fonctionnalités de traceroute et de ping pour fournir une vue complète de la latence, de la perte de paquets et de l'itinéraire des données vers une destination donnée.

![docs/GoTrace.png](docs/GoTrace.png)

## Fonctionnalités 🛠️

Le programme est structuré en trois phases distinctes pour une analyse détaillée :

### Phase 1: Découverte des routeurs (Traceroute)

Cette phase identifie tous les routeurs intermédiaires (les "sauts" ou "hops") entre votre machine et la destination finale. Elle utilise des paquets ICMP avec un Time-to-Live (TTL) incrémentiel pour cartographier l'itinéraire complet.

### Phase 2: Mesure des performances (Ping)

Une fois l'itinéraire tracé, le programme envoie un nombre spécifié de paquets ICMP à chaque routeur de la liste. Il collecte des données de latence pour chaque saut, ce qui permet d'identifier les points de faiblesse ou les goulots d'étranglement sur le chemin.

### Phase 3: Affichage des statistiques

Les résultats sont présentés dans un tableau compact et facile à lire. Pour chaque routeur, l'outil affiche :

L'adresse IP du routeur.

La latence moyenne en microsecondes (µs).

L'écart-type de la latence en microsecondes (µs), qui indique la variabilité des performances.

Le pourcentage de perte de paquets, un indicateur clé de la fiabilité de la connexion.

## Prérequis 📋

- Go (version 1.18 ou supérieure)

- Privilèges d'administrateur (sudo sur Linux/macOS) ou équivalent sur d'autres systèmes, car le programme nécessite l'accès à des sockets ICMP bruts.

## Installation et Utilisation 🚀

Cloner le dépôt :

```Bash
git clone https://github.com/votre_utilisateur/go-network-analyzer.git
```

Lancer le programme :
Exécutez l'application en spécifiant l'hôte cible (nom de domaine ou adresse IP) et le nombre de paquets à envoyer.

```Bash
sudo go run main.go <hostname_ou_ip> <nombre_de_paquets> [délai_en_ms]
```

Exemple :
Pour analyser le chemin vers google.com en envoyant 10 paquets à chaque saut, utilisez la commande suivante :

```Bash
sudo go run main.go google.com 10
```

## Exemple de sortie 📊

Voici à quoi ressemble le résultat final de l'analyse :

```Bash
   ----- Tracing routers to www.google.com (142.250.178.132) ----- 
Hop   | IP Address       | Avg (µs)     | Std Dev (µs)    | Loss (%)  
---------------------------------------------------------------------
1     | 192.168.1.1      | 528.25       | 104.78          | 0.00
2     | 10.12.0.1        | 1256.74      | 258.91          | 0.00
3     | 172.16.25.1      | 2530.12      | 450.32          | 0.00
...
10    | 142.250.75.14    | 25687.55     | 1205.80         | 0.00
```
