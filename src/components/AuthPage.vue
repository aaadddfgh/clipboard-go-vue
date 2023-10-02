<script setup lang="ts">
import { getKey, getAuth } from '@/api';
import { handshakeAESDecrypt, transpotAES, handshakePasswordEncrypt, RSA } from '@/lib/crypto';
import { useRSAStore } from '@/stores/key';
import { onMounted, onBeforeMount, ref } from 'vue';
import { useRouter } from 'vue-router';


const loading = ref(false)
const router = useRouter()
const password = ref("")
const needPass = ref(true)
// {
//     router.push( {name:"clipboard"})
//     //router.push("clipboard")
// }
const ServerRSAPub = useRSAStore();
onBeforeMount(async () => {
    const data = (await getKey())
    ServerRSAPub.key = data.key;
    needPass.value = data.password;
    if (!data.password) {

        alert("没有密码，数据可能会不安全")
        auth();
    }
})


function authProcess(data: { ok: boolean, key: string }) {
    if (data.ok) {
        const aes = handshakeAESDecrypt(data.key)
        transpotAES.key = aes.key;
        transpotAES.iv = aes.iv;
        router.push({ 'name': 'clipboard' })
    }
    else {
        alert("出错了")
        location.reload()
    }
}

function auth(pass?: string) {
    loading.value = true;
    const MyRsa = RSA;
    try {
        if (pass && needPass.value) {
            getAuth(
                {
                    pass: handshakePasswordEncrypt(pass, ServerRSAPub.key),
                    pubKey: MyRsa.getPubKey()
                },
                authProcess
            )
        }
        else {
            getAuth(
                {
                    pubKey: MyRsa.getPubKey()
                },
                authProcess
            )
        }
    }
    catch {

    }
}
</script>

<template>
    <div class="d-flex text-center m-auto" v-if="loading">
        <div>
            <div class="spinner-border m-auto" role="status">
                <span class="sr-only">Loading...</span>
            </div>
            <h4 class="mx-auto">载入中</h4>
        </div>
    </div>

    <div class="card mx-auto" v-else-if="needPass">


        <div class="card-body">
            <h4 class="card-title text-center"> 请输入密码</h4>

            <input type="password" v-model="password" />
            <button v-on:click="auth(password)" type="button" class="mr-3 mt-3 btn btn-primary btn-lg btn-block">接入</button>

        </div>


    </div>
    <div class="d-flex text-center m-auto" v-else>
        <div>
            <div class="spinner-border m-auto" role="status">
                <span class="sr-only">Loading...</span>
            </div>
            <h4 class="mx-auto">载入中</h4>
        </div>
    </div>
</template>

<style scoped></style>