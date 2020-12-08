package main

import (
  "fmt" 
  "os"
//  "image"
  "image/jpeg"
  "log"
)

type Pixel struct{
  r uint32
  g uint32
  b uint32
  a uint32
}

func main() {
  argument := os.Args[1]
  catFile, err := os.Open(argument)
  if err != nil {
	  log.Fatal(err) // ça renvoie une erreur si image introuvable
  }
  defer catFile.Close()

  imData , err := jpeg.Decode(catFile)
  if err != nil {
        fmt.Println(err)
  }
    fmt.Println(imData.At(5,5))
    fmt.Println(imData.Bounds())

  c := imData.At(5,5)
  r,g,b,a := c.RGBA()
  fmt.Printf("c:%v,r:%v\n",c,r/255)
  

  fmt.Println(r/256,g,b,a)
  valeur := uint8(int(r))
  fmt.Println(valeur)
  fmt.Println(r)

  /*taille := imData.Bounds()
	hauteur := taille.Dy()
	largeur := taille.Dx()
  fmt.Println(hauteur)
  fmt.Println(largeur)
  
  for i:=0; i<largeur; i++ {
		for j:=0; j<hauteur; j++{
			 imData.At(i,j).RGBA())
		}
	}*/

   /* var ValeurPixel [256]int
    ValeurPixel[0] = imData.At(0,0)
    fmt.Println(ValeurPixel[0])*/

    //cat, err := jpeg.Decode(catFile)
    //if err != nil {
      //  log.Fatal(err)
    //}
    //fmt.Println(cat)
}