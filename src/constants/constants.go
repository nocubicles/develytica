package constants

func GetAdDirections() map[uint]string {
	return map[uint]string{
		1: "Ost",
		2: "Müük",
		3: "Vahetus",
		4: "Annan ära",
	}
}

func GetAdLocations() map[uint]string {
	return map[uint]string{
		1:  "Harju maakond",
		2:  "Tartu maakond",
		3:  "Ida-Viru maakond",
		4:  "Pärnu maakond",
		5:  "Lääne-Viru maakond",
		6:  "Viljandi maakond",
		7:  "Rapla maakond",
		8:  "Võru maakond",
		9:  "Saare maakond",
		10: "Jõgeva maakond",
		11: "Järva maakond",
		12: "Valga maakond",
		13: "Põlva maakond",
		14: "Lääne maakond",
		15: "Hiiu maakond",
		16: "Muu",
	}
}

func GetAdTypes() map[uint]string {
	return map[uint]string{
		1: "Jalgrattad",
		2: "Jalgratta lisavarustus",
		3: "Jalgratta riided",
		4: "Jalgratta varuosad",
		5: "Muu"}
}

func GetBikeTypes() map[uint]string {
	return map[uint]string{
		1: "Maantee",
		2: "Maastiku",
		3: "Elektri",
		4: "Bmx",
		5: "Hübriid",
		6: "Linna",
		7: "Laste",
		8: "Muu"}
}

func GetFrameSizes() map[uint]string {
	return map[uint]string{
		1: "XXS",
		2: "XS",
		3: "S",
		4: "M",
		5: "L",
		6: "XL",
		7: "XXL",
		8: "Muu"}
}

func GetFrameMaterials() map[uint]string {
	return map[uint]string{
		1: "Carbon",
		2: "Alu",
		3: "Teras",
		4: "Titaan",
		5: "Muu",
	}
}
