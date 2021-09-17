package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("examjwt123")

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func runCommand(params string) string {
	cmd := exec.Command("bash", "-c", params)
	stdout, _ := cmd.Output()

	//if err != nil {
	//	fmt.Println("Ocurrio un error")
	//	fmt.Println(err.Error())
	//	//log.Fatal(err)
	//}

	//fmt.Println(string(stdout))
	//fmt.Println(strings.TrimSuffix(string(stdout), "\n"))
	return string(strings.TrimSuffix(string(stdout), "\n"))

}

func check_run() int {

	runCommand("sudo mkdir -p /var/log/.test")

	if fileExists("/var/log/.test/try1.lock") {
		fmt.Println("Ejecutando por segunda vez la correcion, ya no le queda intentos.")

		runCommand("sudo rm -f /var/log/.test/try1.lock")
		runCommand("sudo touch /var/log/.test/try2.lock")
		return 2

	} else if fileExists("/var/log/.test/try2.lock") {
		fmt.Println("Ya no puede ejecutar la correccion de su examen.")
		os.Exit(0)

	} else {
		fmt.Println("Ejecutando por primera vez la correcion, le queda un intento mas de correccion.")
		runCommand("sudo touch /var/log/.test/try1.lock")
		return 1
	}
	return 0
}

func check_service(points int) int {

	total := 0
	fmt.Println("### Ejercicio 2 ###")
	fmt.Println("")
	stdout1 := runCommand("sudo systemctl is-active firewalld")
	stdout2 := runCommand("sudo systemctl is-enabled firewalld")
	if stdout1 == "inactive" && stdout2 == "disabled" {
		total = total + points
		fmt.Printf("\tFirewalld bien configurado: %d puntos\n", points)

	} else {
		fmt.Printf("\tFirewalld no fue configurado: %d puntos\n", 0)
	}

	return total
}

func check_swap(points int) int {
	total := 0
	fmt.Println("### Ejercicio 3 ###")
	fmt.Println("")
	stdout1 := runCommand("sudo swapon -s | grep /dev/sda3 | awk '{ print $3 }'")
	if stdout1 != "" && stdout1 == "1048572" {
		total += points
		fmt.Printf("\tSWAP ampliado y activado: %d puntos\n", points)
	} else {
		fmt.Printf("\tSWAP no fue ampliado o activado: %d puntos\n", 0)
	}

	return total

}

func check_groups(points int) int {
	total := 0
	fmt.Println("### Ejercicio 6 ###")
	fmt.Println("")
	group1 := runCommand("cat /etc/group | grep -v '#' | grep contabilidad | awk -F ':' '{ print $1 }'")
	group2 := runCommand("cat /etc/group | grep -v '#' | grep monitoreo | awk -F ':' '{ print $1 }'")

	if group1 == "contabilidad" && group2 == "monitoreo" {
		total += points
		fmt.Printf("\tGrupos creados correctamente: %d puntos\n", points)
	} else {
		fmt.Printf("\tTotalidad de grupos no fueron creados: %d puntos\n", 0)
	}

	return total

}

func check_users(points int) int {
	total := 0
	fmt.Println("### Ejercicio 7 ###")
	fmt.Println("")
	user_name := strings.TrimSpace(runCommand("sudo lslogins -u jperez | grep Username | awk -F ':' '{ print $2}'"))
	user_home := strings.TrimSpace(runCommand("sudo lslogins -u jperez | grep 'Home directory' | awk -F ':' '{ print $2}'"))
	user_shell := strings.TrimSpace(runCommand("sudo lslogins -u jperez | grep 'Shell' | awk -F ':' '{ print $2}'"))
	user_groups := strings.TrimSpace(runCommand("sudo lslogins -u jperez | grep 'Supplementary groups' | awk -F ':' '{ print $2}'"))

	if user_name == "jperez" && user_home == "/home/jperez" && user_shell == "/bin/bash" && user_groups == "contabilidad" {
		total += points / 2
		fmt.Printf("\tUsuario jperez creado correctamente: %d puntos\n", points/2)
	} else {
		fmt.Printf("\tUsuario jperez NO creado correctamente: %d puntos\n", 0)
	}

	user_name = strings.TrimSpace(runCommand("sudo lslogins -u nagios | grep Username | awk -F ':' '{ print $2}'"))
	user_home = strings.TrimSpace(runCommand("sudo lslogins -u nagios | grep 'Home directory' | awk -F ':' '{ print $2}'"))
	user_shell = strings.TrimSpace(runCommand("sudo lslogins -u nagios | grep 'Shell' | awk -F ':' '{ print $2}'"))
	user_groups = strings.TrimSpace(runCommand("sudo lslogins -u nagios | grep 'Supplementary groups' | awk -F ':' '{ print $2}'"))

	if user_name == "nagios" && user_home == "/raid/nagios" && user_shell == "/sbin/nologin" && user_groups == "monitoreo" {
		total += points / 2
		fmt.Printf("\tUsuario nagios creado correctamente: %d puntos\n", points/2)
	} else {
		fmt.Printf("\tUsuario nagios NO creado correctamente: %d puntos\n", 0)
	}

	return total

}

func check_account(points int) int {
	total := 0
	fmt.Println("### Ejercicio 8 ###")
	fmt.Println("")
	user_name := strings.TrimSpace(runCommand("sudo lslogins -u jperez | grep Username | awk -F ':' '{ print $2}'"))
	pass_min := strings.TrimSpace(runCommand("sudo lslogins -u jperez | grep 'Minimum change time' | awk -F ':' '{ print $2}'"))
	pass_max := strings.TrimSpace(runCommand("sudo lslogins -u jperez | grep 'Maximum change time' | awk -F ':' '{ print $2}'"))
	pass_warn := strings.TrimSpace(runCommand("sudo lslogins -u jperez | grep 'Password expiration warn interval' | awk -F ':' '{ print $2}'"))
	if user_name == "jperez" && pass_min == "10" && pass_max == "60" && pass_warn == "5" {
		total += points
		fmt.Printf("\tLa cuenta jperez fue configurado correctamente: %d puntos\n", points)
	} else {
		fmt.Printf("\tLa cuenta jperez no fue configurado correctamente: %d puntos\n", 0)
	}
	return total
}

func check_kernel_params(points int) int {
	total := 0
	fmt.Println("### Ejercicio 5 ###")
	fmt.Println("")
	stdout1 := runCommand("sysctl --values net.ipv4.ip_forward")
	stdout2 := runCommand("sysctl --values kernel.sysrq")
	if stdout1 == "1" && stdout2 == "1" {
		total += points / 2
		fmt.Printf("\tParametros del kernel establecidos correctamente: %d puntos\n", points/2)
	} else {
		fmt.Printf("\tParametros del kernel no establecidos: %d puntos\n", 0)
	}

	stdout1 = runCommand("sudo cat /etc/sysctl.conf | grep 'net.ipv4.ip_forward'")
	stdout2 = runCommand("sudo cat /etc/sysctl.conf | grep 'kernel.sysrq'")

	if stdout1 != "" && stdout2 != "" {
		total += points / 2
		fmt.Printf("\tParametros del kernel establecidos de forma persistente correctamente: %d puntos\n", points/2)
	} else {
		fmt.Printf("\tParametros del kernel no establecidos de forma persistente correctamente: %d puntos\n", 0)
	}

	return total

}

func check_permissions(points int) int {
	total := 0
	fmt.Println("### Ejercicio 9 ###")
	fmt.Println("")

	stdout := runCommand("sudo ls -ld /raid/nagios")
	if stdout != "" {
		perms := strings.Fields(stdout)[0]
		owner := strings.Fields(stdout)[2]
		group := strings.Fields(stdout)[3]
		if perms == "drwxrwxr--." && owner == "nagios" && group == "monitoreo" {
			total += points / 2
			fmt.Printf("\tDirectorio /raid/nagios configurado correctamente: %d puntos\n", points/2)
		} else {
			fmt.Printf("\tDirectorio /raid/nagios NO configurado correctamente: %d puntos\n", 0)
		}
	} else {
		fmt.Printf("\tDirectorio /raid/nagios NO creado, resolver el ejecicio 3: %d puntos\n", 0)
	}

	stdout = runCommand("sudo ls -ld /raid/contab")
	if stdout != "" {
		perms := strings.Fields(stdout)[0]
		owner := strings.Fields(stdout)[2]
		group := strings.Fields(stdout)[3]
		if perms == "drwxrws---." && owner == "root" && group == "contabilidad" {
			total += points / 2
			fmt.Printf("\tDirectorio /raid/contab configurado correctamente: %d puntos\n", points/2)
		} else {
			fmt.Printf("\tDirectorio /raid/contab NO configurado correctamente: %d puntos\n", 0)
		}
	} else {
		fmt.Printf("\tDirectorio /raid/contab NO creado, resolver el ejecicio 3: %d puntos\n", 0)
	}

	return total
}

func check_localrepo(points int) int {
	total := 0
	fmt.Println("### Ejercicio 11 ###")
	fmt.Println("")
	stdout := runCommand("yum repolist --enabled | grep localrepo")
	if stdout != "" {
		total += points
		fmt.Printf("\tRepositorio localrepo configurado correctamente: %d puntos\n", points)
	} else {
		fmt.Printf("\tRepositorio localrepo NO configurado correctamente: %d puntos\n", 0)
	}

	return total
}

func check_disablerepo(points int) int {
	total := 0
	fmt.Println("### Ejercicio 12 ###")
	fmt.Println("")
	stdout := runCommand("yum repolist --enabled | grep extras")
	if stdout == "" {
		total += points
		fmt.Printf("\tRepositorio extras deshabilitado correctamente: %d puntos\n", points)
	} else {
		fmt.Printf("\tRepositorio extras NO deshabilitado correctamente: %d puntos\n", 0)
	}

	return total
}

func check_initramfs(points int) int {
	total := 0
	fmt.Println("### Ejercicio 13 ###")
	fmt.Println("")

	initramfs := runCommand("sudo ls /boot/initramfs-$(uname -r).img")
	initramfs_bkp := initramfs + ".bkp"

	if fileExists(initramfs_bkp) {
		total += points / 2
		fmt.Printf("\tBackup de initramfs generado correctamente: %d puntos\n", points/2)
	} else {
		fmt.Printf("\tBackup de initramfs NO generado correctamente: %d puntos\n", 0)
	}

	new_initramfs := runCommand("sudo ls -l /boot/initramfs-$(uname -r).img")
	system_month := runCommand("sudo date +%b")
	system_day := runCommand("sudo date +%d")
	newinitramfs_month := strings.Fields(new_initramfs)[5]
	newinitramfs_day := strings.Fields(new_initramfs)[6]

	if system_month == newinitramfs_month && system_day == newinitramfs_day {
		total += points / 2
		fmt.Printf("\tNuevo initramfs generado correctamente: %d puntos\n", points/2)
	} else {
		fmt.Printf("\tNuevo initramfs NO fue generado correctamente: %d puntos\n", 0)
	}

	return total

}

func check_httprepo(points int) int {
	total := 0
	fmt.Println("### Ejercicio 10 ###")
	fmt.Println("")
	stdout := runCommand("sudo rpm -qa | grep httpd")
	if stdout == "" {
		fmt.Printf("\tServidor Web NO instalado: %d puntos\n", 0)
	} else {

		// httpd configuration
		stdout1 := runCommand("sudo systemctl is-active httpd")
		stdout2 := runCommand("sudo systemctl is-enabled httpd")
		if stdout1 == "active" && stdout2 == "enabled" {
			total += points / 2
			fmt.Printf("\tServidor Web configurado exitosamente: %d puntos\n", points/2)
		} else {
			fmt.Printf("\tServidor Web NO configurado exitosamente: %d puntos\n", 0)
		}

		// web repo configuration
		stdout1 = runCommand("ls -ld /var/www/html/yum/CentOS/el7/base/repodata/")
		if stdout1 != "" {
			total += points / 2
			fmt.Printf("\tRepositorio creado correctamente: %d puntos\n", points/2)
		} else {
			fmt.Printf("\tRepositorio NO creado correctamente: %d puntos\n", 0)
		}

	}

	return total
}

func check_raid(points int) int {
	total := 0
	fmt.Println("### Ejercicio 4 ###")
	fmt.Println("")

	// md general
	stdout1 := runCommand("sudo cat /proc/mdstat | grep md")
	if stdout1 != "" {

		// md0 raid level
		stdout_raidlevel := runCommand("cat /proc/mdstat | grep md0 | awk '{ print $4 }'")
		if stdout_raidlevel == "raid1" {
			total += points / 6
			fmt.Printf("\t/dev/md0 fue configurado con el raid correcto: %d puntos\n", points/6)
		} else {
			fmt.Printf("\t/dev/md0 no fue configurado con el raid correcto: %d puntos\n", 0)
		}

		// md0 mounted
		stdout_mounted := runCommand("df -hT | grep /dev/md0 | grep -v Filesystem")
		if stdout_mounted != "" {
			raid_name := strings.Fields(stdout_mounted)[0]
			raid_fs := strings.Fields(stdout_mounted)[1]
			raid_mount := strings.Fields(stdout_mounted)[6]
			if raid_name == "/dev/md0" && raid_fs == "xfs" && raid_mount == "/raid/contab" {
				total += points / 6
				fmt.Printf("\t/dev/md0 fue configurado con el punto de montaje correcto: %d puntos\n", points/6)
			} else {
				fmt.Printf("\t/dev/md0 no fue configurado con el punto de montaje correcto: %d puntos\n", 0)
			}
		} else {
			fmt.Printf("\t/dev/md0 no fue montado: %d puntos\n", 0)
		}

		// md0 fstab
		stdout_fstab := runCommand("cat /etc/fstab | grep /raid/contab")
		if stdout_fstab != "" {
			raid_uuid := strings.Fields(stdout_fstab)[0]
			if strings.HasPrefix(raid_uuid, "UUID=") {
				total += points / 6
				fmt.Printf("\t/dev/md0 fue configurado persistentemente con su UUID: %d puntos\n", points/6)
			} else {
				fmt.Printf("\t/dev/md0 NO fue configurado persistentemente con su UUID de forma correcta: %d puntos\n", 0)
			}
		} else {
			fmt.Printf("\t/dev/md0 NO fue configurado persistentemente de forma correcta: %d puntos\n", 0)
		}

		// md1 raid level
		stdout_raidlevel = runCommand("cat /proc/mdstat | grep md1 | awk '{ print $4 }'")
		if stdout_raidlevel == "raid5" {
			total += points / 6
			fmt.Printf("\t/dev/md1 fue configurado con el raid correcto: %d puntos\n", points/6)
		} else {
			fmt.Printf("\t/dev/md1 no fue configurado con el raid correcto: %d puntos\n", 0)
		}

		// md1 mounted
		stdout_mounted = runCommand("df -hT | grep /dev/md1 | grep -v Filesystem")
		if stdout_mounted != "" {
			raid_name := strings.Fields(stdout_mounted)[0]
			raid_fs := strings.Fields(stdout_mounted)[1]
			raid_mount := strings.Fields(stdout_mounted)[6]
			if raid_name == "/dev/md1" && raid_fs == "xfs" && raid_mount == "/raid/nagios" {
				total += points / 6
				fmt.Printf("\t/dev/md1 fue configurado con el punto de montaje correcto: %d puntos\n", points/6)
			} else {
				fmt.Printf("\t/dev/md1 no fue configurado con el punto de montaje correcto: %d puntos\n", 0)
			}
		} else {
			fmt.Printf("\t/dev/md1 no fue montado: %d puntos\n", 0)
		}

		// md1 fstab
		stdout_fstab = runCommand("cat /etc/fstab | grep /raid/nagios")
		if stdout_fstab != "" {
			raid_uuid := strings.Fields(stdout_fstab)[0]
			if strings.HasPrefix(raid_uuid, "UUID=") {
				total += points / 6
				fmt.Printf("\t/dev/md1 fue configurado persistentemente con su UUID: %d puntos\n", points/6)
			} else {
				fmt.Printf("\t/dev/md1 NO fue configurado persistentemente con su UUID de forma correcta: %d puntos\n", 0)
			}
		} else {
			fmt.Printf("\t/dev/md1 NO fue configurado persistentemente de forma correcta: %d puntos\n", 0)
		}

	} else {
		fmt.Printf("\tRAID no fue configurado: %d puntos\n", 0)
	}

	return total
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	//claims := token.Claims.(jwt.MapClaims)

	//claims["authorized"] = true
	//claims["client"] = "Elliot Forbes"
	//claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func post(try int, score int, student string) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	fmt.Println("Enviando resultado al cloud")
	todo := Todo{strconv.Itoa(try), strconv.Itoa(score), student}

	jsonReq, err := json.Marshal(todo)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://ferztyle.rocks:8080/api/exam", bytes.NewBuffer(jsonReq))
	fmt.Println(validToken)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf(string(body))
}

// Todo struct
type Todo struct {
	Try     string `json:"try"`
	Score   string `json:"score"`
	Student string `json:"student"`
}

func read_student_name() string {
	if fileExists("/etc/student") {
		name := runCommand("sudo cat /etc/student | tr -d '\n'")
		if name == "" {
			fmt.Printf("\n\n##### DEBE COMPLETAR PRIMERO EL ARCHIVO /etc/student CON SU NOMBRE Y APELLIDO #####\n\n")
			log.Fatal()
		} else {
			return name
		}
	} else {
		fmt.Printf("\n\n##### DEBE COMPLETAR PRIMERO EL ARCHIVO /etc/student CON SU NOMBRE Y APELLIDO #####\n\n")
		log.Fatal()
	}
	return ""
}

func main() {

	student := read_student_name()

	try := check_run()
	// 1- root password recovered
	fmt.Println("### Ejercicio 1 ###")
	fmt.Println("")
	total := 2
	fmt.Printf("\tPassword root recuperado: %d puntos\n\n", total)
	// 2- Services: Stopping and disabling
	total += check_service(5)
	fmt.Println("")
	// 3- Increase Swap area
	total += check_swap(5)
	fmt.Println("")
	// 4- Configure Raid
	total += check_raid(24)
	fmt.Println("")
	// 5- Kernel Params settings
	total += check_kernel_params(10)
	fmt.Println("")
	// 6- Add Groups
	total += check_groups(5)
	fmt.Println("")
	// 7- Add users
	total += check_users(10)
	fmt.Println("")
	// 8- User password policies
	total += check_account(5)
	fmt.Println("")
	// 9- Directories permissions
	total += check_permissions(10)
	fmt.Println("")
	// 10- Publish local repo via apache
	total += check_httprepo(10)
	fmt.Println("")
	// 11- Set Local repo
	total += check_localrepo(5)
	fmt.Println("")
	// 12- Disable repo
	total += check_disablerepo(5)
	fmt.Println("")
	// 13- Recreate initramfs
	total += check_initramfs(4)
	fmt.Println("")
	// Total
	fmt.Printf("Total obtenido: %d\n\n", total)

	post(try, total, student)
}
