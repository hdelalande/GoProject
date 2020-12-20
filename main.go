package main

import (
  "fmt"
  "image"
  "image/color"
  "os"
//  "image"
  "image/jpeg"
  "log"
  //"math"
  //"sync"
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

  var testNoirEtBlanc bool
  taille := imData.Bounds()
	hauteur := taille.Dy()
  largeur := taille.Dx()
  //restelargeur := math.Mod(float64(largeur),100) 
  testNoirEtBlanc = true
  //NBboucles := (largeur/100)+1
  //var wg sync.WaitGroup
  //wg.Add(NBboucles)
  for i:=0; i<largeur; i++ {
    for j:=0; j<hauteur; j++{
      r,g,b,a := imData.At(i,j).RGBA()
      if r!=g || g!=b || r!=b { /*test si les intensités RGB sont différente pour detecter si l'image est en couleur*/
        testNoirEtBlanc = false 
        fmt.Println("Tu dois nous envoyer une image en noir et blanc bg et l'opacité est de", a)
        break
      } else {
        continue
      }
    }
    break
  }

  if testNoirEtBlanc == false {
    /*Dans cette partie, nous transformons l'image en noir est blanc si ce n'est pas le cas*/
    imgSet := image.NewRGBA(taille) //on commence par créer une image "vide" de la même taille que l'image d'origine.
    for y := 0; y < taille.Max.Y; y++ { // deux boucles pour parcourir l'ensemble des pixels constituant l'image.
      for x := 0; x < taille.Max.X; x++ {
        oldPixel := imData.At(x, y) // on récupère le pixel à la position x,y.
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
    testNoirEtBlanc = true
  }
  //wg.Wait()
}
func histogramme(Data image.Image) [65536]uint32{
    var ValuePixel [65536]uint32 // Tableau qui va permettre de savoir combien il y aura de pixels pour chaque intensité
    taille := Data.Bounds()
	  hauteur := taille.Dy()
    largeur := taille.Dx()

    // Dans cette boucle, on compte le nbr de pixels par intensité
    for i := 0; i < largeur; i++ {
      for j := 0; j < hauteur; j++ {
        r, _, _, _ := Data.At(i, j).RGBA() // On récupère la valeur du pixel en RGBA
        ValuePixel[r] = ValuePixel[r] + 1 // On ajoute 1 à l'index d'intensité r
      }
    }
    return ValuePixel
  }
func ProbaPixel(ValuePixel []uint32, NombrePixel int) [65536]float32{
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

func normalisation(ProbaPixelCumul []float32, Data image.Image) [65536]float32{
    var ImageNorma [65536]float32 // Tableau contenant les intensités de pixels normalisés
    // ImageNorma[z] = x
    // z correspond à l'intensité du pixel sur l'image de base
    // x sera la nouvelle intensité pour l'image normalisée

    // Dans cette boucle, on calcule les nouvelles inensités (égalisé) avec la formule
    taille := Data.Bounds()
	  hauteur := taille.Dy()
    largeur := taille.Dx()
    for i := 0; i < largeur; i++ {
      for j := 0; j < hauteur; j++ {
        r, _, _, _ := Data.At(i, j).RGBA() // Pareil que ma boucle précédente
        ImageNorma[r] = 65535 * ProbaPixelCumul[r] // Formule pour normalisé une image
      }
    }
  return ImageNorma
  }
func creationimage(Data image.Image, ImageNorma []float32)  {
  // Création de l'image normalisée
  taille := Data.Bounds()
  hauteur := taille.Dy()
  largeur := taille.Dx()
  imgSet := image.NewRGBA(taille)
  for i := 0; i < largeur; i++ {
    for j := 0; j < hauteur; j++ {
      r, g, b, _ := Data.At(i, j).RGBA()
      lum := 0.299*float64(ImageNorma[r]) + 0.587*float64(g) + 0.114*float64(b)
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
