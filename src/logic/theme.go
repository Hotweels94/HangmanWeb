// Programmed by:
// - LIENARD Mathieu
// - DORGES Guillaume
// - AMSELLEM--BOUSIGNAC Ryan

package hangman

// Func to find an image in the folder depending on which theme it is and the health points remaining
func GetImagePath(health int, theme int) string {

	// Ducky theme
	if theme == 2 {
		switch health {
		case 10:
			return "./img/theme2/d1.png"
		case 9:
			return "./img/theme2/d2.png"
		case 8:
			return "./img/theme2/d3.png"
		case 7:
			return "./img/theme2/d4.png"
		case 6:
			return "./img/theme2/d5.png"
		case 5:
			return "./img/theme2/d6.png"
		case 4:
			return "./img/theme2/d7.png"
		case 3:
			return "./img/theme2/d8.png"
		case 2:
			return "./img/theme2/d9.png"
		case 1:
			return "./img/theme2/d9.png"
		case 0:
			return "./img/theme2/d10.png"
		default:
			return "./img/theme2/d1.png"
		}
	}

	// Halloween theme
	if theme == 3 {
		switch health {
		case 10:
			return "./img/theme3/h11.png"
		case 9:
			return "./img/theme3/h10.png"
		case 8:
			return "./img/theme3/h9.png"
		case 7:
			return "./img/theme3/h8.png"
		case 6:
			return "./img/theme3/h7.png"
		case 5:
			return "./img/theme3/h6.png"
		case 4:
			return "./img/theme3/h5.png"
		case 3:
			return "./img/theme3/h4.png"
		case 2:
			return "./img/theme3/h3.png"
		case 1:
			return "./img/theme3/h2.png"
		case 0:
			return "./img/theme3/h1.png"
		default:
			return "./img/theme3/h1.png"
		}
	}

	// Christmas theme
	if theme == 4 {
		switch health {
		case 10:
			return "./img/theme4/n11.png"
		case 9:
			return "./img/theme4/n10.png"
		case 8:
			return "./img/theme4/n9.png"
		case 7:
			return "./img/theme4/n8.png"
		case 6:
			return "./img/theme4/n7.png"
		case 5:
			return "./img/theme4/n6.png"
		case 4:
			return "./img/theme4/n5.png"
		case 3:
			return "./img/theme4/n4.png"
		case 2:
			return "./img/theme4/n3.png"
		case 1:
			return "./img/theme4/n2.png"
		case 0:
			return "./img/theme4/n1.png"
		default:
			return "./img/theme4/n1.png"
		}
	}
	return "./img/theme2/d1.png"
}
