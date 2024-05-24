package run

import (
	"fmt"
	"io/ioutil"
	"monitor/service"
	"os"
	"strings"
)

func Run(config string) {

	// Delete file "info.txt" if it exists
	err1 := os.Remove("info.txt")
	if err1 != nil && !os.IsNotExist(err1) {
		fmt.Println("Error deleting file:", err1)
		return
	}

	WentWrong := map[string]bool{
		"srv":  false,
		"web":  false,
		"disk": false,
	}

	var Data []string

	// Service checker
	srv, _ := service.NewServiceChecker()

	srv.LoadConfig(config)
	srv.ResetNotRunning()

	err := srv.CheckStatus()
	if err != nil {
		fmt.Println(err)
	}

	srvslice := srv.NotRunning
	isEmpty := len(srvslice) == 0
	if !isEmpty {
		WentWrong["srv"] = true
		// Add the line "following services not running:" and a new line break to the Data slice
		Data = append(Data, "following services not running:")
		Data = append(Data, "") // Empty line

		// Add the elements from srvslice that are not running to Data
		for _, service := range srvslice {
			Data = append(Data, service)

		}
		Data = append(Data, "-----------------------------------------------------------------------------------")
	}

	// WebChecker to check webapp is reachable and cert is valid
	web, _ := service.NewWebChecker()

	web.LoadConfig(config)
	web.ResetWebErrors()

	web.CheckWebsite()
	web.CheckCert()

	webSlice := web.WebErrors
	isEmpty = len(webSlice) == 0
	if !isEmpty {
		WentWrong["web"] = true
		// Add the line "following services not running:" and a new line break to the Data slice
		Data = append(Data, "following WebApp errors occurred:")
		Data = append(Data, "") // Empty line

		// Add the elements from srvslice that are not running to Data
		for _, problem := range webSlice {
			Data = append(Data, problem)

		}
		Data = append(Data, "-----------------------------------------------------------------------------------")
	}

	// Check Disk space
	disk, _ := service.NewSysChecker()

	disk.LoadConfig(config)
	disk.ResetLowDiskSpace()

	disk.CheckDisks()

	diskSlice := disk.LowDiskSpace
	isEmpty = len(diskSlice) == 0
	if !isEmpty {
		WentWrong["disk"] = true
		// Add the line "following services not running:" and a new line break to the Data slice
		Data = append(Data, "Diskspace low on following Disks:")
		Data = append(Data, "") // Empty line

		// Add the elements from srvslice that are not running to Data
		for _, dsk := range diskSlice {
			Data = append(Data, dsk)

		}
		Data = append(Data, "-----------------------------------------------------------------------------------")
	}

	fmt.Println(Data)

	// Write the Data slice to the file "info.txt"
	err = ioutil.WriteFile("info.txt", []byte(strings.Join(Data, "\r\n")), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Println("Data has been written to the file info.txt.")

	// Check if anything is set to true in the WentWrong slice
	sendMail := false
	for _, v := range WentWrong {
		if v {
			sendMail = true
			break
		}
	}

	// If sendMail is true, then send an email
	if sendMail {
		// Load mail configuration
		mail, err := service.NewMail()
		if err != nil {
			fmt.Println("Error loading mail config:", err)
			return
		}
		mail.LoadConfig(config)

		// Create email message
		subject := "Monitor found some Problems"
		body := "Check attachment info.txt"

		// Read the content of the file "info.txt"
		attachment, err := ioutil.ReadFile("info.txt")
		if err != nil {
			fmt.Println("Error reading attachment:", err)
			return
		}

		// Send email with attachment
		err = mail.SendMailWithAttachment([]byte(body), attachment, subject)
		if err != nil {
			fmt.Println("Error sending email:", err)
			return
		}

		fmt.Println("Email has been sent.")
	}
}
