package main

import (
  "fmt"
  "image"
  "image/color"
  "os"
//  "image"
  "image/jpeg"
  "log"
  "sync"
)


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

  var testNoirEtBlanc bool
  var tabIntens [65536]uint32
  taille := imData.Bounds()
	hauteur := taille.Dy()
  largeur := taille.Dx() 
  NBboucles := hauteur*largeur
  testNoirEtBlanc = true
  var wg sync.WaitGroup
  wg.Add(NBboucles)
  for i:=0; i<largeur; i++ {
    for j:=0; j<hauteur; j++{
      r,g,b,a := imData.At(i,j).RGBA()
      if r!=g || g!=b || r!=b { /*test si les intensités RGB sont différente pour detecter si l'image est en couleur*/
        testNoirEtBlanc = false 
        fmt.Println("Tu dois nous envoyer une image en noir et blanc bg")
        break
      } else {
        continue
      }
      break
    }
    wg.Wait()
  }
  if testNoirEtBlanc == false {
    /*Dans cette partie on peut imaginer un code permettant de faire passer une image en N&B avant de faire le traitement de l'image*/
    imgSet := image.NewRGBA(taille)
    for y := 0; y < taille.Max.Y; y++ {
      for x := 0; x < taille.Max.X; x++ {
        oldPixel := imData.At(x, y)
        r, g, b, _ := oldPixel.RGBA()
        lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
        pixel:= color.Gray{uint8(lum / 256)}
        imgSet.Set(x, y, pixel)
      }

    }

    outFile, err := os.Create("changed.jpg")
    if err != nil {
      log.Fatal(err)
    }
    defer outFile.Close()
    jpeg.Encode(outFile, imgSet, nil)
    testNoirEtBlanc:= true

  }
  if testNoirEtBlanc == true {
        var ValuePixel [65536]float32

        for i := 0; i < largeur; i++ {
          for j := 0; j < hauteur; j++ {
            r, g, b, a := imData.At(i, j).RGBA()
            salut := g + b + a
            fmt.Println("Yo les gars : ", salut)
            ValuePixel[r] = ValuePixel[r] + 1
            fmt.Println("Valeur du tab en ", r, " : ", ValuePixel[r])
          }
        }

        var total float32
        for z := 0; z < 65536; z++ {
          total += ValuePixel[z]
        }

        fmt.Println("Nbr total pixel : ", total)
        fmt.Println("Largeur :", largeur, "Hauteur : ", hauteur)

        //var ValuePixelEqua [65536]float32

        /*for i := 0; i < largeur; i++ {
          for j := 0; j < hauteur; j++ {
              ValuePixelEqua[]
          }
        }*/
      }

      //test

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

  