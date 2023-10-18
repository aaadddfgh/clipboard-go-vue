const url= import.meta.env.DEV ? "/api":""

export function getAuth(data: { pass?: string,/* pubKey: string ,*/ aes:string}, action: (data: { ok: boolean, /*key: string*/ }) => any) {
    // WARNING: For POST requests, body is set to null by browsers.

    var xhr = new XMLHttpRequest();


    xhr.addEventListener("readystatechange", function () {
        if (this.readyState === 4) {
            action(JSON.parse(this.responseText));
        }
    });

    xhr.open("POST", url+"/auth");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.send(JSON.stringify(data));
}

export function getContent(action: (data: { data: string }) => any) {
    var xhr = new XMLHttpRequest();


    xhr.addEventListener("readystatechange", function () {
        if (this.readyState === 4) {
            action(JSON.parse(this.responseText));
        }
    });

    xhr.open("GET", url+"/content");
    xhr.send()
}
export function postContent(data: { test: string, content: string }) {
    // WARNING: For POST requests, body is set to null by browsers.
    return new Promise<{ok:boolean}>((res, rej) => {

        var xhr = new XMLHttpRequest();


        xhr.addEventListener("readystatechange", function () {
            if (this.readyState === 4) {
                
                    res(JSON.parse(this.responseText));
                
            }
        });

        xhr.open("POST", url+"/content");
        xhr.setRequestHeader("Content-Type", "application/json");

        xhr.send(JSON.stringify(data));
    })
}

export async function getKey() {
    return await new Promise<{ key: string,password:boolean,lang:string }>((res) => {
        var xhr = new XMLHttpRequest();


        xhr.addEventListener("readystatechange", function () {
            if (this.readyState === 4) {

                res(JSON.parse(this.responseText))
            }
        });

        xhr.open("GET", url+"/key");

        xhr.send();
    })

}

export function postFile(data: FormData) {
    // WARNING: For POST requests, body is set to null by browsers.
    return new Promise<{ok:boolean}>((res, rej) => {

        var xhr = new XMLHttpRequest();


        xhr.addEventListener("readystatechange", function () {
            if (this.readyState === 4) {
                
                    res(JSON.parse(this.responseText));
                
            }
        });

        xhr.open("POST", url+"/file");

        xhr.send(data);
    })
}

export function getFile() {
    window.open(url+"/file")
}

export function delFile() {
    var xhr = new XMLHttpRequest();

    xhr.open("GET", url+"/delfile");

    xhr.send();
}