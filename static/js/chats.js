const url = "http://localhost:8080"
const urlDeleteChat = url + "/chat/delete/"
const urlSaveSettingChat = url + "/chat/setting/save/"
const urlSwapChats = url + "/chat/swap"
const urtSaveTitle = url + "/chat/edit"
let objDiv = document.getElementById("field-chat");
objDiv.scrollTop = objDiv.scrollHeight;
let activeChatID = 0

function deleteChatById(id, element) {
    fetch(urlDeleteChat+id, {method: "DELETE"})
        .then(response => response.json())
        .then(commits => {
            if (commits.status === "ok") {
                element.parentNode.removeChild(element);
            } else {
                console.log(commits.message)
            }
        });
}

function saveTitleForChat(id, title) {
    fetch(urtSaveTitle, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify({
            "id": id,
            "title": title,
        }),
    })
        .then(response => response.json())
        .then(commits => {
            if (commits.status !== "ok") {
                console.log(commits.code, commits.message)
            }
        });
}

function swapChats(down_elem, up_elem) {
    let id1 = down_elem.id.substring(5);
    let id2 = up_elem.id.substring(5);
    console.log("swap", id1, id2)
    fetch(urlSwapChats, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify({
            "id1": id1,
            "id2": id2,
        }),
    })
        .then(response => response.json())
        .then(commits => {
            if (commits.status === "ok") {
                up_elem.parentNode.removeChild(up_elem);
                down_elem.insertAdjacentElement( 'afterend', up_elem );
                add_event_for_swap_button();
            } else {
                console.log(commits.code, commits.message)
            }
        });
}

function saveSettingChatById(id) {
    console.log(speechFlag.checked, volumeFlag.checked, roleFlag.value, parseFloat(temparatureFlag.value))
    fetch(urlSaveSettingChat+id, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify({
            "speech": speechFlag.checked,
            "volume": volumeFlag.checked,
            "role": roleFlag.value,
            "temperature": parseFloat(temparatureFlag.value),
        }),
    })
        .then(response => response.json())
        .then(commits => {
            if (commits.status === "error") {
                console.log(commits.code, commits.message)
            }
        });
}

document.querySelectorAll(".chat").forEach((element) => {
    let state = true;
    if (element.classList.contains("active")) {
        activeChatID = element.id
        element.getElementsByClassName("delete-chat")[0].style.visibility = "hidden";
    } else {
        element.onclick = (event) => {location.href="/main/"+element.id.substring(5)}
        element.getElementsByClassName("delete-chat")[0].addEventListener("click", (event)=>{
            event.stopPropagation();
            deleteChatById(element.id.substring(5), element);
        })
    }
    let img = element.getElementsByClassName("edit-img")[0];
    let mid = element.getElementsByClassName("mid-title")[0];
    element.getElementsByClassName("edit-chat")[0].addEventListener("click", (event)=>{
        event.stopPropagation();
        if (state) {
            state = false;
            mid.removeChild(mid.children[0])
            mid.insertAdjacentHTML("beforeend", `<input type="text" style="width:100%">`);
            mid.getElementsByTagName("input")[0].addEventListener("click", (event)=>{
                event.stopPropagation();
            })
            img.src="/static/img/check.svg";
        } else {
            state = true;
            let title = mid.getElementsByTagName("input")[0].value;
            mid.removeChild(mid.children[0]);
            mid.insertAdjacentHTML("beforeend", `<p>${title}</p>`);
            saveTitleForChat(element.id.substring(5), title);
            img.src="/static/img/pencil-square.svg";
        }
    })
});

function add_event_for_swap_button() {
    let chats = document.querySelectorAll(".chat")
    chats.forEach((element, ind) => {
        if (ind !== 0) {
            element.getElementsByClassName("up-chat")[0].style.visibility = "visible";
            element.getElementsByClassName("up-chat")[0].onclick = (event)=>{
                event.stopPropagation();
                swapChats(element, chats[ind - 1]);
            }
        } else {
            element.getElementsByClassName("up-chat")[0].style.visibility = "hidden";
        }
        if (ind !== chats.length - 1) {
            element.getElementsByClassName("down-chat")[0].style.visibility = "visible";
            element.getElementsByClassName("down-chat")[0].onclick = (event)=>{
                event.stopPropagation();
                swapChats(chats[ind + 1], element);
            }
        } else {
            element.getElementsByClassName("down-chat")[0].style.visibility = "hidden";
        }
    });
}
add_event_for_swap_button();

const speechFlag = document.getElementById("myonoffswitch")
const volumeFlag = document.getElementById("myonoffswitch2")
const roleFlag = document.getElementById("state")
const temparatureFlag = document.getElementById("flying")

if (speechFlag.checked) {
    document.getElementById("button-send").children[0].src="/static/img/mic_black.svg";
}

speechFlag.addEventListener("change", (event) => {
    saveSettingChatById(activeChatID.substring(5))
})
volumeFlag.addEventListener("change", (event) => {
    saveSettingChatById(activeChatID.substring(5))
})
roleFlag.addEventListener("change", (event) => {
    saveSettingChatById(activeChatID.substring(5))
})
temparatureFlag.addEventListener("change", (event) => {
    saveSettingChatById(activeChatID.substring(5))
})
