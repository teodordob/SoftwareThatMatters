# Powershell script to invoke the shell script
param(
    [int]$number
)

if($number){
    wsl -e sh ./init.sh $number
} else {
    wsl -e sh ./init.sh -1
}
