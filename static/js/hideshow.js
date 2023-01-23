function hideShow(elem) {
    console.log("Something happened");
    var x = document.getElementById("user-posts");
    var a = document.getElementById("posts")
    var y = document.getElementById("user-comments");
    var b = document.getElementById("comments")
    var z = document.getElementById("liked-user-posts");
    var c = document.getElementById("likes")
    if (elem.value == "posts") {
        x.style.display = "block";
        y.style.display = "none";
        z.style.display = "none";
        a.style.textDecoration = "underline"
        b.style.textDecoration = "none"
        c.style.textDecoration = "none"
    } else if (elem.value == "comments") {
        x.style.display = "none";
        y.style.display = "block";
        z.style.display = "none";
        a.style.textDecoration = "none"
        b.style.textDecoration = "underline"
        c.style.textDecoration = "none"
    } else if (elem.value == "liked-posts") {
        x.style.display = "none";
        y.style.display = "none";
        z.style.display = "block";
        a.style.textDecoration = "none"
        b.style.textDecoration = "none"
        c.style.textDecoration = "underline"
    }
}