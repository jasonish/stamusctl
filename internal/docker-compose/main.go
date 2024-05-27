package compose

func isValideRestartModeDocker(restart string) bool {
	if restart != "no" &&
		restart != "always" &&
		restart != "on-failure" &&
		restart != "unless-stopped" {
		return false
	}
	return true
}

func isValideRamSizeDocker(ram string) bool {
	return ram != "" && ram != "0" && ram != "0m" && ram != "0g	" && ram != "0k" && ram != "0t" && ram != "0p"
}
