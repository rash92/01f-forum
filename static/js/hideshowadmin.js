
function hideShowAdmin(elem) {
    console.log("Hello world")
    var x = document.getElementById("mod-requests");
    var y = document.getElementById("user-admin");
    var z = document.getElementById("tags-admin");
    var a = document.getElementById("mod")
    var b = document.getElementById("users")
    var c = document.getElementById("tags")
    if (elem.value == "mod") {
        x.style.display = "block";
        y.style.display = "none";
        z.style.display = "none";
        a.style.textDecoration = "underline"
        b.style.textDecoration = "none"
        c.style.textDecoration = "none"
    } else if (elem.value == "users") {
        x.style.display = "none";
        y.style.display = "block";
        z.style.display = "none";
        a.style.textDecoration = "none"
        b.style.textDecoration = "underline"
        c.style.textDecoration = "none"
    } else if (elem.value == "tags") {
        x.style.display = "none";
        y.style.display = "none";
        z.style.display = "block";
        a.style.textDecoration = "none"
        b.style.textDecoration = "none"
        c.style.textDecoration = "underline"
    }
}