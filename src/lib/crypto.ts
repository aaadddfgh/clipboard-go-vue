import forge from 'node-forge'

class RSACypto{
    public keypair:forge.pki.rsa.KeyPair;
    constructor(){
        this.keypair = forge.pki.rsa.generateKeyPair({bits: 2048})
    }
    getPubKey(){
        return forge.pki.publicKeyToPem(this.keypair.publicKey)
    }
    encypt(data:string){
        return this.keypair.publicKey.encrypt(data);
    }
    decypt(data:string){
        return this.keypair.privateKey.decrypt(data)
    }
}


export class AESCypto{
    key = forge.random.getBytesSync(16);
    iv = forge.random.getBytesSync(16);


    constructor(/*key:string,iv:string*/){

    }
    
    getKeyForTransport(pubKey: string) {
        return EncryptDataWithRSAPubKey(
            JSON.stringify({ key: forge.util.bytesToHex(this.key), iv: forge.util.bytesToHex(this.iv) }),
            pubKey
        );

    }

    encypt(data:string){
        let cipher = forge.cipher.createCipher('AES-CBC', this.key);
        cipher.start({iv: this.iv});
        cipher.update(forge.util.createBuffer(data));
                cipher.finish();
        ;
        
        // outputs encrypted hex
        
        return forge.util.bytesToHex(cipher.output.data);
    }
    encryptBytes(data:string){
        let cipher = forge.cipher.createCipher('AES-CBC', this.key);
        cipher.start({iv: this.iv});
        cipher.update(forge.util.createBuffer(data));
                cipher.finish();
        ;
        
        // outputs encrypted hex
        
        return cipher.output.data;
    }
    decypt(data:string){
        let decipher = forge.cipher.createDecipher('AES-CBC', this.key);
        decipher.start({iv: this.iv});
        decipher.update(forge.util.createBuffer(forge.util.hexToBytes(data)));
        var result = decipher.finish(); // check 'result' for true/false
        // outputs decrypted hex
        
        return decipher.output.data;
    }
}

export const RSA=new RSACypto()

export function handshakeAESDecrypt(AES:string,):{key:string,iv:string}{
    const data= forge.util.hexToBytes(AES);
    const d:{key:string,iv:string}=JSON.parse(RSA.decypt(data));
  
    return {
        key: forge.util.hexToBytes(d.key),
        iv:forge.util.hexToBytes(d.iv),
    };
}

export function EncryptDataWithRSAPubKey(data:string,pubKey:string){
    return forge.util.bytesToHex(
        forge.pki.publicKeyFromPem(pubKey)
            .encrypt(
                data
            )
    )
}

export function handshakePasswordEncrypt(pass:string,pubKey:string){
    // HEX( RSA( AES(pass)))
    return forge.util.bytesToHex(
        forge.pki.publicKeyFromPem(pubKey)
            .encrypt(
                //注意 password被aes
                transpotAES.encryptBytes(pass)
            )
        );
}

export const transpotAES=new AESCypto();