package main

import (
    "fmt"
    "image"
    "image/color"
    "os"
  //"image"
    "image/jpeg"
    "log"
    "sync"
    "math"
)

func main() {
  argument := os.Args[1]
  catFile, err := os.Open(argument)
  if err != nil {
    log.Fatal(err) /* Renvoi une erreur si il n'y a pas de fichier à ouvrir en argument */
  }
  defer catFile.Close()

  imData, err := jpeg.Decode(catFile)
  if err != nil {
    fmt.Println(err) /* Renvoi une erreur si le fichier n'est pas une image jpeg ou jpg */
  }


  taille := imData.Bounds()
  hauteur := taille.Dy()
  largeur := taille.Dx()
  NbPixel := largeur*hauteur
  Decoupe := 100
  r := make([][]uint32,NbPixel)
  g := make([][]uint32,NbPixel)
  b := make([][]uint32,NbPixel)
  a := make([][]uint32,NbPixel)
  for i:=0; i<largeur; i++ {
    for j:=0; j<hauteur; j++{
    r[i][j],g[i][j],b[i][j],a[i][j] = imData.At(i,j).RGBA()
      if r[i][j]!=g[i][j] || g[i][j]!=b[i][j] || r[i][j]!=b[i][j] { /*test si les intensités RGB sont différente pour detecter si l'image est en couleur*/
        fmt.Println("Tu dois nous envoyer une image en noir et blanc bg et l'opacité est de", a)
        imData = passagenoiretblanc(imData)
        break
      } 
    }
  }
  var NBboucles int
  restelargeur := math.Mod(float64(largeur),float64(Decoupe))
  if restelargeur != 0{
    NBboucles = (largeur/Decoupe)+1
  } 
  if restelargeur == 0{
    NBboucles = (largeur/Decoupe)
  }
  var HistogrammePixel [65536]uint32
  //var TabImageEga [65536]float32
  var ImageEgalise [65536]float32
  //c1 := make(chan [65536]uint32)
  //c2 := make(chan [65536]float32)
  var wg1 sync.WaitGroup
  wg1.Add(NBboucles)
  for i:=largeur ; i>0; i=i-Decoupe{
    if i>Decoupe{
      go  histogramme(&HistogrammePixel,r,&hauteur,(i-Decoupe),i,&wg1)
    }
    if i<Decoupe{
      go histogramme(&HistogrammePixel,r,&hauteur,0,i,&wg1)
    }
  }
  wg1.Wait()
  TabDesProba := probapixel(HistogrammePixel,NbPixel)
  var wg2 sync.WaitGroup
  wg2.Add(NBboucles)
  for i:=largeur ; i>0; i=i-Decoupe{
    if i>Decoupe{
      go egalisation(&ImageEgalise,&TabDesProba,r,&hauteur,(i-Decoupe),i,&wg2)
    }
    if i<Decoupe{
      go egalisation(&ImageEgalise,&TabDesProba,r,&hauteur,0,i,&wg2)
    }
  }
  /*
  var tabC2 [][65536]float32
  for b:=0; b<NBboucles; b++{
    tabC2[b] = <- c2
    for p:=0; p<65536; p++{
      TabImageEga[p] += tabC2[b][p]
    }
  }
  */
  creationimage(imData,ImageEgalise)
}


func histogramme(ValuePixel *[65536]uint32,tab [][]uint32, hauteur *int, largeur1 int, largeur2 int, wg1 *sync.WaitGroup){
 // Tableau qui va permettre de savoir combien il y aura de pixels pour chaque intensité
  // Dans cette boucle, on compte le nbr de pixels par intensité
  for i := largeur1; i < largeur2; i++ {
    for j := 0; j < *hauteur; j++ {
      r:=tab[i][j] // On récupère la valeur du pixel en RGBA
        ValuePixel[r] = ValuePixel[r] + 1
    }
  }
  wg1.Done()
  // c1 <- ValuePixel
}

func probapixel(ValuePixel [65536]uint32, NombrePixel int) [65536]float32{
  var ProbaPixelCumul [65536]float32 // Tableau des probas cumulées
  var max uint32 // Variable permettant de faire la cumulation des probas
  max = 0 
  // Dans cette boucle, on calcul la probabilité cumulée de chaque pixel
  for i := 0; i < 65536; i++ {
    max = max + ValuePixel[i] // Pour faire la proba cumulée
    ProbaPixelCumul[i] = float32(max) / float32(NombrePixel) // On divise ce nbr de pixels par le nombre total de pixel
  }
  return ProbaPixelCumul
}

func egalisation(ImageEga *[65536]float32, ProbaPixelCumul *[65536]float32,tab [][]uint32,hauteur *int,largeur1 int, largeur2 int, wg2 *sync.WaitGroup){
  // Tableau contenant les intensités de pixels égalisés
  // ImagEga[z] = x
  // z correspond à l'intensité du pixel sur l'image de base
  // x sera la nouvelle intensité pour l'image égalisée

  // Dans cette boucle, on calcule les nouvelles inensités (égalisées) avec la formule
  for i := largeur1; i < largeur2; i++ {
    for j := 0; j < *hauteur; j++ {
      r:=tab[i][j] // On récupère la valeur du pixel en RGBA
      ImageEga[r] = 65535 * ProbaPixelCumul[r] // Formule pour normalisé une image
    }
  }
  wg2.Done()
  //c2 <- ImageEga
}

func creationimage(Data image.Image, ImageEga [65536]float32)  {
  // Création de l'image égalisée
  taille := Data.Bounds()
  hauteur := taille.Dy()
  largeur := taille.Dx()
  imgSet := image.NewRGBA(taille)
  for i := 0; i < largeur; i++ {
    for j := 0; j < hauteur; j++ {
      r, g, b, _ := Data.At(i, j).RGBA()
      lum := 0.299*float64(ImageEga[r]) + 0.587*float64(g) + 0.114*float64(b)
      pixel := color.Gray16{uint16(lum)}
      imgSet.Set(i,j, pixel)
    }
  }

  outFile, err := os.Create("test.jpg")
  if err != nil {
    log.Fatal(err)
  }
  defer outFile.Close()
  jpeg.Encode(outFile, imgSet, nil)

}

func passagenoiretblanc(Data image.Image) image.Image{
  /*Dans cette partie, nous transformons l'image en noir est blanc si ce n'est pas le cas*/
  taille := Data.Bounds()
  imgSet := image.NewRGBA(taille) //on commence par créer une image "vide" de la même taille que l'image d'origine.
  for y := 0; y < taille.Max.Y; y++ { // deux boucles pour parcourir l'ensemble des pixels constituant l'image.
    for x := 0; x < taille.Max.X; x++ {
      oldPixel := Data.At(x, y) // on récupère le pixel à la position x,y.
      r, g, b, _ := oldPixel.RGBA() // on recupère les valeurs d'intensité en rouge, vert et bleu.
      lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b) // on calcule l'intensité la mieux adaptée grace à une formule.
      pixel:= color.Gray16{uint16(lum)} // on fait appel à color.Gray pour transformer le pixel en gris
      imgSet.Set(x, y, pixel) // 
    }
  }
  outFile, err := os.Create("changed.jpg")
  if err != nil {
    log.Fatal(err)
  }
  defer outFile.Close()
  jpeg.Encode(outFile, imgSet, nil)
  return imgSet
}

    
    /* //////////// QUELQUES LIGNES POUR DEBUGUER EN CAS DE PROBLEME ////////////

    var total uint32
    for z := 0; z < 65536; z++ {
      total += ValuePixel[z]
    }

    for m := 0 ; m < 65356 ; m++ {
      if ImageNorma[m] != 0 {
        fmt.Println("En ", m, " ca vaut ", ImageNorma[m])
      }
    }

    var cumul uint32
    cumul = 0
    for o := 0 ; o < 39835 ; o++ {
      cumul = cumul + ValuePixel[o]
    }
    fmt.Println("Nbr pixel avant ::::", cumul)
    fmt.Println("ValuePixel ::::", ValuePixel[39835])
    fmt.Println("ImageNorma ::::", ImageNorma[39835])
    fmt.Println("Cumulated : ", ValuePixelEqua[39835])

    fmt.Println("Nbr total pixel : ", total)
    fmt.Println("Largeur :", largeur, "Hauteur : ", hauteur)
    fmt.Println("Maxi : ", max)

    fmt.Println(hauteur)
    fmt.Println(largeur)
    fmt.Println(imData.At(5,5))
    fmt.Println(imData.Bounds())

    var ValeurPixel [256]int
    ValeurPixel[0] = imData.At(0,0)
    fmt.Println(ValeurPixel[0])

    cat, err := jpeg.Decode(catFile)
    if err != nil {
    log.Fatal(err)
    }
    fmt.Println(cat)

    //////////// QUELQUES LIGNES POUR DEBUGUER EN CAS DE PROBLEME //////////// */