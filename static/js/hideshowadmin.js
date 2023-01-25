
function hideShowAdmin(elem) {
    var x = document.getElementById("mod-requests");
    var a = document.getElementById("mod")
    var y = document.getElementById("user-admin");
    var b = document.getElementById("users")
    var z = document.getElementById("tag-admin");
    var c = document.getElementById("tags")
    if (elem.value == "mod-requests") {
        x.style.display = "block";
        y.style.display = "none";
        z.style.display = "none";
        a.style.textDecoration = "underline"
        b.style.textDecoration = "none"
        c.style.textDecoration = "none"
    } else if (elem.value == "cser-admin") {
        x.style.display = "none";
        y.style.display = "block";
        z.style.display = "none";
        a.style.textDecoration = "none"
        b.style.textDecoration = "underline"
        c.style.textDecoration = "none"
    } else if (elem.value == "tag-admin") {
        x.style.display = "none";
        y.style.display = "none";
        z.style.display = "block";
        a.style.textDecoration = "none"
        b.style.textDecoration = "none"
        c.style.textDecoration = "underline"
    }
}