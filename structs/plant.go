package structs

// Plant contains information about a single plant
type Plant struct {
	ID         int    `json:"id" binding:"required"`
	CommonName string `json:"common-name" binding:"required"`
	// Likes       int      `json:"likes"`
	// Images      []string `json:"images"`
	// Description string   `json:"description"`
	// SpeciesName string   `json:"species-name"`
	// FamilyName  string   `json:"family-name"`
	// Collections []string `json:"collections"`
}

// Kingdom
//   -> Subkingdom
//     -> Division
//       -> Division class
//         -> Division order
//           -> Family
//             -> Genus
//               -> Plant

// For example, the balsam fir hierarchy is:

// Kingdom -> Plantae – (Plants)
// Subkingdom -> Tracheobionta – (Vascular plants)
// Division -> Coniferophyta – (Conifers)
// Class -> Pinopsida
// Order -> Pinales
// Family -> Pinaceae – (Pine family)
// Genus -> Abies
// Plant -> Abies balsamea

// We'll create a list of plants
// var plants = []Plant{
// 	Plant{1, 0, "", ""},
// 	Plant{2, 0, "", ""},
// 	Plant{3, 0, "", ""},
// 	Plant{4, 0, "", ""},
// 	Plant{5, 0, "", ""},
// 	Plant{6, 0, "", ""},
// 	Plant{7, 0, "", ""},
// }
