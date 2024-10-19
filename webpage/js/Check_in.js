//簽到退切換
let ischeckin = true;
document.getElementById("Switch_mode").addEventListener('click', function(){
    const switchbtn = document.getElementById("Switch_mode");
    const submitbtn = document.getElementById("Check_in_Btn");

    if(ischeckin){
        switchbtn.innerHTML = "切換為簽到";
        submitbtn.innerHTML = "簽退";
        ischeckin = false;
        alert("切換為簽退");
    }else{
        switchbtn.innerHTML = "切換為簽退";
        submitbtn.innerHTML = "簽到";
        ischeckin = true;
        alert("切換為簽到");
    }
})
