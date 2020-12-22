# GoProject
Ce programme traite des images de type JPEG.
L'image passe par un test pour connaître sa nature (B&W ou en couleur). Si elle est en couleur, l'image est transformée en B&W. Ensuite, l'image est égalisée (On ajuste le contraste de l'image à l'aide d'une formule). 
Pour exécuter le programme :
go run <nom du fichier.go> <nom de l'image .jpeg>
Vous aurez en sortie une ou deux nouvelles images :
- si l'image est initialement en noir et blanc, vous aurez une nouvelle image intitulée changed.jpg égalisée.
- si l'image est en couleur, vous aurez deux images en sorties : une après passage en noir et blanc, et une deuxième égalisée.
  
# Compostion du GitHub

Vous trouverez dans le GitHub deux programmes :

- mainsansgoroutine.go : le programme fonctionne sans goroutine
- main.go : le programme fonctionne avec des goroutines

Les fichiers restants sont des images en B&W ou couleur afin de pouvoir les mettre en entré du programme.
