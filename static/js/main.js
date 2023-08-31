import { CustomWebSocket } from "/static/js/custom-web-socket.js";
window.MediaRecorder = OpusMediaRecorder;
const workerOptions = {
    OggOpusEncoderWasmPath: 'https://cdn.jsdelivr.net/npm/opus-media-recorder@0.8.0/OggOpusEncoder.wasm',
    WebMOpusEncoderWasmPath: 'https://cdn.jsdelivr.net/npm/opus-media-recorder@0.8.0/WebMOpusEncoder.wasm'
};

const btnQuestion = document.getElementById("button-send");
const btnStop = document.getElementById("button-send");
const consoleWithQuestionText = document.getElementById("floatingTextarea2")
const consoleWithAnswerText = document.getElementById("answer-text");
let ws = new CustomWebSocket();
let questionText = "";
let isRecordNow = false;
let isQuestionReceived = false;
let player = document.querySelector('#player');
let voice = [];
let globMediaRecorder = null;
const speechFlag = document.getElementById("myonoffswitch")
let btnState = 1;
if (speechFlag.checked) {
    btnState = 0;
}
let haveMedia = false
speechFlag.addEventListener("change", (event) => {
    if (!speechFlag.checked) {
        if (btnState === 0) {
            btnState = 1;
            btnQuestion.children[0].src="/static/img/send-alt-1.svg";
        }
        if (isRecordNow) {
            isRecordNow = false;
            globMediaRecorder.stop();
        }
    }
    if (speechFlag.checked && !isRecordNow) {
        consoleWithQuestionText.value="";
        btnState = 0;
        btnQuestion.children[0].src="/static/img/mic_black.svg";
    }
})

consoleWithQuestionText.addEventListener("click", (ev)=>{
    if (isRecordNow) {
        isRecordNow = false;
        globMediaRecorder.stop();
    } else {
        if (btnState === 0) {
            btnState = 1;
            btnQuestion.children[0].src="/static/img/send-alt-1.svg";
        }
    }
})


navigator.mediaDevices.getUserMedia({audio: true, video: false})
    .then(stream => {
        let options = { mimeType: "audio/wave" };
        const mediaRecorder = new MediaRecorder(stream, options, workerOptions);
        globMediaRecorder = mediaRecorder;
        haveMedia = true;

        mediaRecorder.addEventListener("start", (ev) => {
            console.log("FFFFF")
            voice = []
        })

        mediaRecorder.addEventListener("dataavailable", (ev) => {
            // voice.push(ev.data)
            ws.socket.send(ev.data);
            console.log(`chunk send: size ${ev.data.size}`);
            if (mediaRecorder.state === "inactive") {
                console.log("inactive");
                ws.socket.send("DONE");
            }
        });

        // btnStop.addEventListener("click", (ev) => {
        //     console.log("record stop")
        //     mediaRecorder.stop();
        //     console.log(mediaRecorder.mimeType)
        // });
    });


btnQuestion.addEventListener("click", (ev) => {
    console.log(btnState, isRecordNow)
    if (btnState === 0 && haveMedia) {
        console.log("record start");
        btnState = 1;
        btnQuestion.children[0].src="/static/img/send-alt-1.svg";
        questionText = "";
        isRecordNow = true;
        consoleWithQuestionText.value = "";
        // consoleWithAnswerText.innerHTML = "";
        isQuestionReceived = false;
        ws.newConnect("ws://localhost:8080/asr", ()=>{console.log("Open new connect");},
            ()=>{console.log("Close connect");}, (mes) => {
                consoleWithQuestionText.value=mes.data;
                // console.log(mes)
                // console.log(typeof(mes.data));
                if (isQuestionReceived) {
                    if (typeof mes.data === "string") {
                        consoleWithAnswerText.innerHTML += mes.data;
                    } else {
                        voice.push(mes.data)
                        console.log(typeof(voice))
                        let blob = new Blob(voice, {'type': "audio/wave"});
                        voice = []
                        let audioURL = URL.createObjectURL(blob);
                        player.src = audioURL;
                        player.play()
                        player.onended = () => {ws.socket.send(0)}
                        // player.addEventListener("stop")
                        // (await async () => {
                            // await function playAudio(){
                            //     return new Promise( async res=>{
                            //       player.play()
                            //       player.onended = res
                            //     })
                            //   }()
                        // })();
                        console.log("e");
                        // const newLocal = n = new Audio(audioURL);
                        // n.play();
                    }
                } else {
                    let res = JSON.parse(mes.data)
                    // console.log(res)
                    consoleWithQuestionText.value = questionText + res.text;
                    if (res.fix) {
                        questionText += res.text;
                    }
                    if (res.finish) {
                        consoleWithQuestionText.insertAdjacentHTML("beforeend", `<h4>Ожидайте ответ на ваш вопрос</h4><p>${questionText}</p>`);
                        isQuestionReceived = true;
                    }
                }
                // consoleWithAnswerText.innerHTML += mes.data;
            }, (err)=>{alert(err);})
        globMediaRecorder.start(100);
    } else if (btnState === 0) {
            btnState = 1;
            btnQuestion.children[0].src="/static/img/send-alt-1.svg";
            alert("Проверьте микрофон. Сейчас возможна отправка вопроса только текстом.")
    } else if (btnState === 1) {
        if (isRecordNow) {
            console.log("DDDD")
            isRecordNow = false;
            globMediaRecorder.stop();
            console.log("DDDD")
        }
        let question = consoleWithQuestionText.value;
        if (question === "") {
            return
        }
        alert("Запрос отправлен: " + question)
        consoleWithQuestionText.value = "";
        btnState = 2;
        btnQuestion.children[0].src="/static/img/x-square.svg";

    } else {
        alert("Запрос отменён")
        consoleWithQuestionText.value = "";
        if (speechFlag.checked) {
            btnState = 0;
            btnQuestion.children[0].src="/static/img/mic_black.svg";
        } else {
            btnState = 1;
            btnQuestion.children[0].src="/static/img/send-alt-1.svg";
        }
    }
    // console.log("record start");
    // if (isRecordNow) {
    //     mediaRecorder.stop();
    //     isRecordNow = false;
    //     return
    // }
});
