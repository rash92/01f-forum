function hideShow(elem) {
    console.log("Something happened");
    var x = document.getElementById("user-posts");
    var y = document.getElementById("user-comments");
    var z = document.getElementById("liked-user-posts");
    if (elem.value == "posts") {
        x.style.display = "block";
        y.style.display = "none";
        z.style.display = "none";
    } else if (elem.value == "comments") {
        x.style.display = "none";
        y.style.display = "block";
        z.style.display = "none";
    } else if (elem.value == "liked-posts") {
        x.style.display = "none";
        y.style.display = "none";
        z.style.display = "block";
    }
}