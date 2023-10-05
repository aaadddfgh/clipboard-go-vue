<script setup lang="ts">
import { postContent, getContent } from '@/api';
import { transpotAES } from '@/lib/crypto';
import { useRSAStore } from '@/stores/key';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import MessageVue from './Message.vue';
import Quill  from "quill";
const router = useRouter()
//@ts-ignore
var quillInstance ;

if (useRSAStore().key === "") {
    // router.push({name:"auth"})
}

function postBoard(content: string) {
    try {
        postContent(
            {
                test: transpotAES.encypt("test"),
                content: transpotAES.encypt(content)
            },
        ).then(({ ok }) => {
            if (ok) {
                printContent();
                info("放入成功");
            }
            else {
                alert("出错了")
                router.push({ name: "auth" })
            }
        })
    }
    catch {
        alert("出错了")
        router.push({ name: "auth" })
    }
}

//const content = ref("")
let editDiv: HTMLDivElement | null = null;

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
            catch {
                if (confirm("似乎出错了，您要重新进入吗？")) {
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
        placeholder: '向剪贴板放入内容',
        theme: 'snow',  // or 'bubble'
        scrollingContainer:'#editableDiv'
    });
    //window.qqqq=quillInstance;
    editDiv = document.querySelector("#editableDiv")
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
    //@ts-ignore
    (quillInstance as Quill).setContents(JSON.parse(v))
}

function copyContent() {
    
    //@ts-ignore
    (quillInstance as Quill).setSelection(0,(quillInstance as Quill).getLength())
    document.execCommand("copy")
    //navigator.clipboard.writeText(getClipboard());
    info("已复制")
}
</script>

<template>
    <div class="container d-flex flex-column mx-2 p-2 main">
        <Transition name="message">
            <MessageVue v-if="showInfo">
                {{ infoData }}
            </MessageVue>
        </Transition>
        <span class="h4 mb-3">向剪贴板中放入内容</span>
        <div class="shadow text-area mb-3 p-0 border border-dark rounded" id="editableDiv">
            <div id="editor-container"></div>
        </div>
        <div class="mx-auto d-flex">
            <button type="button" class="btn btn-success mx-3" v-on:click="copyContent">复制</button>
            <button type="button" class="btn btn-primary mx-3" v-on:click="//@ts-ignore
                postBoard(getClipboard())">放入内容</button>
            <button type="button" class="btn btn-primary mx-3"
                v-on:click="() => { printContent(() => { info('已刷新') }); }">刷新</button>
        </div>

    </div>
</template>

<style scoped>
.main {
    max-width: 100vw;
    min-width: 70vw;
    min-height: 70vh;
}

.text-area {
    max-width: 80vw;
    overflow: auto;
}

#editor-container{
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
