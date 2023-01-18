function hideShow() {
    console.log("Something happened")
    var x = document.getElementById("user-posts");
    var y = document.getElementById("user-comments")
    if (x.style.display === "none" && y.style.display === "block") {
        x.style.display = "block";
        y.style.display = "none";
    } else {
        x.style.display = "none";
        y.style.display = "block";
    }
}