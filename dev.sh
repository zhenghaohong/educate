set user root
set host 42.138.169.141
set pwd ZhongAn0755.
set fromPath mqtt
set toPath /home/educate/mqtt

spawn scp $fromPath $user@$host:$toPath
expect {
	"*password:" { send "$pwd\n"}

}

interact



