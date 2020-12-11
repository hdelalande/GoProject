package main

import (
  "fmt"
  "image"
  "image/color"
  "os"
//  "image"
  "image/jpeg"
  "log"
)

type pixel struct{ /*Notre structure d'un pixel = son intensité, son opacité, sa position en x et sa position en y*/
  intensPix uint32
  opac uint32
  Dx int
  Dy int 
}

func main() {
  argument := os.Args[1]
  catFile, err := os.Open(argument)
  if err != nil {
	  log.Fatal(err) /* Renvoi une erreur si il n'y a pas de fichier à ouvrir en argument */ 
  }
  defer catFile.Close()

  imData , err := jpeg.Decode(catFile)
  if err != nil {
        fmt.Println(err) /* Renvoi une erreur si le fichier n'est pas une image jpeg ou jpg */
  }

  var list []pixel
  var testNoirEtBlanc bool
  var tabIntens [2^32-1]uint32
  taille := imData.Bounds()
	hauteur := taille.Dy()
  largeur := taille.Dx() 
  /* 
  ici il faudra décomposer la hauteur et la largeur dans d'autres
  variables pour faire  fonctionner les goroutines 
  */
  testNoirEtBlanc = true
  for i:=0; i<largeur; i++ {
    for j:=0; j<hauteur; j++{
      r,g,b,a := imData.At(i,j).RGBA()
      if r!=g || g!=b || r!=b { /*test si les intensités RGB sont différente pour detecter si l'image est en couleur*/
        testNoirEtBlanc = false 
        fmt.Println("Tu dois nous envoyer une image en noir et blanc bg")
        break
      } 
      InfoPixel := pixel{intensPix:r, Dx:i, Dy:j, opac:a} 
      list = append(list, InfoPixel) /*Stockage des infos de chaque pixels dans une liste*/
      tabIntens [r] += 1 /* Génére un tableau ayant comme indice l'insité du pixel et en valeur le nombre total de pixels ayant cette intensité*/
    }
  }
  if testNoirEtBlanc == false {
    /*Dans cette partie on peut imaginer un code permettant de faire passer une image en N&B avant de faire le traitement de l'image*/
    imgSet := image.NewRGBA(taille)
    for y := 0; y < taille.Max.Y; y++ {
      for x := 0; x < taille.Max.X; x++ {
        oldPixel := imData.At(x, y)
        r, g, b, _ := oldPixel.RGBA()
        y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
        pixel := color.Gray{uint8(y / 256)}
        imgSet.Set(x, int(y), pixel)
      }
    }

    outFile, err := os.Create("changed.jpg")
    if err != nil {
      log.Fatal(err)
    }
    defer outFile.Close()
    jpeg.Encode(outFile, imgSet, nil)

  }
    if testNoirEtBlanc == true {
    /*histogramme, égalisation et normalisation des pixels*/
    }
    fmt.Println(list[9].intensPix)
    fmt.Println(testNoirEtBlanc)
}

//test

  /*
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
    */

  