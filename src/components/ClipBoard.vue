<script setup lang="ts">
import { postContent, getContent, getFile, postFile as poFile, delFile } from '@/api';
import { transpotAES } from '@/lib/crypto';
import { useRSAStore } from '@/stores/key';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import MessageVue from './Message.vue';
import Quill from "quill";
const router = useRouter()
//@ts-ignore
var quillInstance;
import { useI18n } from 'vue-i18n'
import type { MessageSchema } from '@/i18n/type';

const showFile = ref(false)

const { t } = useI18n<{
    message: MessageSchema,
}>()
if (useRSAStore().key === "") {
    // router.push({name:"auth"})
}
var uploadFile: File | null;
function postBoard(content: string) {
    try {
        //使用encodeURI以避免有中文时会出错
        postContent(
            {
                test: transpotAES.encypt("test"),
                content: transpotAES.encypt(encodeURI(content))
            },
        ).then(({ ok }) => {
            if (ok) {
                printContent();
                info(t("z放入成功信息"));
            }
            else {
                alert(t("z错误提示"))
                router.push({ name: "auth" })
            }
        })
    }
    catch {
        alert(t("z错误提示"))
        router.push({ name: "auth" })
    }
}
function fileChange(payload: Event) {
    const fileList = (payload.target as any).files as FileList;
    const file = fileList.item(0);
    if (file) {
        uploadFile = file;
    }
    else {
        uploadFile = null;
    }

}
//const content = ref("")
function postFile() {
    const data = new FormData()
    if (uploadFile !== null) {
        data.append("file", uploadFile)
        data.append("key", transpotAES.encypt("check"))
    }
    poFile(data)
}

function printContent(action?: (recivedData?: string) => void) {

    getContent(
        (data) => {
            try {
                //content.value = transpotAES.decypt(data.data);
                //@ts-ignore
                setClipboaed(transpotAES.decypt(data.data));
                if (action) {
                    //action(content.value)
                    //@ts-ignore
                    action(/*editDiv.innerHTML*/)
                }
            }
            catch  {
                if (confirm(t("z出错重进提示"))) {
                    router.push({ 'name': 'auth' })
                }
            }
        })

}
onMounted(() => {
    quillInstance = new Quill('#editor-container', {
        modules: {
            toolbar: [
                ['bold', 'italic', 'underline', 'strike'],        // toggled buttons
                ['blockquote', 'code-block', "image"],

                [{ 'header': 1 }, { 'header': 2 }],               // custom button values
                [{ 'list': 'ordered' }, { 'list': 'bullet' }],
                [{ 'script': 'sub' }, { 'script': 'super' }],      // superscript/subscript
                [{ 'indent': '-1' }, { 'indent': '+1' }],          // outdent/indent
                [{ 'direction': 'rtl' }],                         // text direction

                [{ 'size': ['small', false, 'large', 'huge'] }],  // custom dropdown
                [],

                [{ 'color': [] }, { 'background': [] }],          // dropdown with defaults from theme
                [{ 'font': [] }],
                [{ 'align': [] }],

                ['clean']                                         // remove formatting button
            ]
        },
        placeholder: t('z输入占位'),
        theme: 'snow',  // or 'bubble'
        scrollingContainer: '#editableDiv'
    });
    //window.qqqq=quillInstance;

    printContent();
})


const showInfo = ref(false)
const infoData = ref("")
function info(msg: string) {
    showInfo.value = false;
    showInfo.value = true;
    infoData.value = msg;
    setTimeout(() => {
        showInfo.value = false
    }, 3000);
}

function getClipboard(): string {
    //@ts-ignore
    return JSON.stringify((quillInstance as Quill).getContents());
}

function setClipboaed(v: string) {
    try{
         
    //@ts-ignore 
    (quillInstance as Quill).setContents(JSON.parse(decodeURI(v)))
    }
    catch(err){
        console.log(err)
    }
}

function copyContent() {

    //@ts-ignore
    (quillInstance as Quill).setSelection(0, (quillInstance as Quill).getLength())
    document.execCommand("copy")
    //navigator.clipboard.writeText(getClipboard());
    info(t("z复制后信息"))
}


const ws = new WebSocket("ws://" + location.host + "/ws");
// 连接成功
ws.onopen = () => {
    console.log('WebSocket连接成功')
}
// 接收到消息
ws.onmessage = (e) => {
    printContent()
    info(t('z刷新后信息'))
};

function goTop(){document.querySelector('#top')!.scrollIntoView();}
function goBottom(){document.querySelector('#bottom')!.scrollIntoView();}
</script>

<template>
    <div class="container d-flex flex-column mx-2 p-2 main">
        <Transition name="message">
            <MessageVue v-if="showInfo">
                {{ infoData }}
            </MessageVue>
        </Transition>
        <span id="top" class="h4 mb-3">{{ t('z输入占位') }}</span>
        <div class="shadow text-area mb-3 p-0 border border-dark rounded mx-auto" id="editableDiv">
            <div id="editor-container"></div>
        </div>
        <div id="bottom" class="mx-auto d-flex flex-wrap">
            <button type="button" class="btn btn-success my-1 mx-3" v-on:click="copyContent">{{ t("z复制按钮") }}</button>
            <button type="button" class="btn btn-primary my-1 mx-3" v-on:click="//@ts-ignore
                postBoard(getClipboard())">{{ t("z放入剪贴板按钮") }}</button>
            <button type="button" class="btn btn-primary my-1 mx-3"
                v-on:click="() => { printContent(() => { info(t('z刷新后信息')) }); }">{{ t("z刷新按钮") }}</button>
            <button type="button" class="btn btn-success my-1 mx-3" v-on:click="showFile = true">放入文件</button>
        </div>

        <div class="affix-area" >
            <div v-on:click="goTop" class="affix bg-secondary rounded-circle d-flex justify-content align-middle mb-1" href="#">
                <svg xmlns="http://www.w3.org/2000/svg" width="auto" height="auto" fill="white" class="bi bi-arrow-up" viewBox="0 0 16 16">
                    <path fill-rule="evenodd" d="M8 15a.5.5 0 0 0 .5-.5V2.707l3.146 3.147a.5.5 0 0 0 .708-.708l-4-4a.5.5 0 0 0-.708 0l-4 4a.5.5 0 1 0 .708.708L7.5 2.707V14.5a.5.5 0 0 0 .5.5z"/>
                  </svg>
                </div>

            <div v-on:click="goBottom" class="affix bg-secondary rounded-circle d-flex justify-content align-middle mb-1">
                <svg xmlns="http://www.w3.org/2000/svg" width="auto" height="auto" fill="white" class="bi bi-arrow-down" viewBox="0 0 16 16">
                    <path fill-rule="evenodd" d="M8 1a.5.5 0 0 1 .5.5v11.793l3.146-3.147a.5.5 0 0 1 .708.708l-4 4a.5.5 0 0 1-.708 0l-4-4a.5.5 0 0 1 .708-.708L7.5 13.293V1.5A.5.5 0 0 1 8 1z"/>
                  </svg>
            </div>
        </div>

        <div class="modal show fade" data-toggle="modal" id="modal" tabindex="-1" v-if="showFile">
            <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable modal-xl" id="dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title" id="staticBackdropLabel">查看文件</h5>
                        <button type="button" class="close" v-on:click="showFile = false">
                            <span>&times;</span>
                        </button>
                    </div>
                    <div class="modal-body d-flex justify-content flex-wrap">
                        <!-- <label for="file">上传文件</label> -->
                        <input id="file" type="file" class="w-100 mt-1" v-on:change="fileChange">
                        <button type="button" class="btn btn-danger mt-1 ml-3" v-on:click="delFile">清空剪贴板中的文件</button>
                    </div>
                    <div class="modal-footer flex-wrap">
                        <button type="button" class="btn btn-primary mx-3" v-on:click="getFile">查看文件</button>
                        <button type="button" class="btn btn-success mx-3" v-on:click="postFile">上传文件</button>
                        <button type="button" class="btn btn-danger mx-3" v-on:click="showFile = false">关闭窗口</button>

                    </div>
                </div>
            </div>
        </div>


    </div>
</template>

<style scoped>
#modal {
    display: inherit;
    background: rgba(0, 0, 0, 0.5)
}

.affix-area{
    position: fixed;  
    bottom: 20px; 
    left:10px
}
.affix{
     width: 30px; 
     height: 30px;
     cursor: pointer;
     padding: 6px;
}
.main {
    max-width: 100vw;
    min-width: 70vw;
    min-height: 70vh;
}

.text-area {
    max-width: 80vw;
    overflow: auto;
}

#editor-container {
    min-height: 40vh;
}

.message-enter-active,
.message-leave-active {
    transition: all 0.5s ease;
}

.message-enter-from,
.message-leave-to {
    transform: translateY(-50px);
    opacity: 20;
}
</style>
