# GoProject
Ce programme traite des images de type JPEG. Il prend en entrée n'importe quelle image (de type jpeg).
L'image passe par un test pour connaître sa nature (Noir et blanc ou en couleurs). Si elle est en couleurs,
l'image est transformée en noir et blanc.
Ensuite, l'image est égalisée (On ajuste le contraste de l'image à l'aide d'une formule). 
Pour exécuter le programme :
go run <nom du fichier.go> <nom de l'image .jpeg>
Vous aurez en sortie une ou deux nouvelles images :
- si l'image est initialement en noir et blanc, vous aurez une nouvelle image intitulé changed.jpg égalisée.
- si l'image est en couleurs, vous aurez deux images en sorties : une après passage en noir et blanc, et une deuxième égalisée.