# gomli

![gomli](assets/gomli.png)

`go + smali ~= gomli` a dynamic smali patching utility

## Usage

A simple replace operation `gomli replace -s script.js -d work/com.example_decompiled`

Or you are welcome to use `--help`

```
A go parser for smali code

Usage:
  gomli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  find        find call chains XREFs for a particular class
  help        Help about any command
  replace     replace the values of const-strings globally in the application

Flags:
  -h, --help      help for gomli
  -v, --verbose   enable verbose/debug mode

Use "gomli [command] --help" for more information about a command.
```

## Replace

Supported instructions: 
* `const-string`
* `fill-array-data`

In order do maintain some code consistency when using a decompiler tool like `jadx` specific rules when it comes to value replacement have to be considered to each instruction, obviously `const-string` is the most simple just replace the string argument from the instruction. Arrays `fill-array-data` are a bit trickier.

Lets say we have the following code
```
    private int A(JSONObject jSONObject) {
      return jSONObject.optInt(y(new byte[]{-118, -94, 76, -22, -26, -50, 41, -119}), 0);
    }

    public static String y(byte[] bArr) {
      int length = bArr.length / 2;
      int length2 = bArr.length - length;
      int i2 = length2 - length;
      byte[] bArr2 = new byte[length];
      int i3 = 0;
      while (i3 < length) {
        bArr2[i3] = (byte) ((bArr[length2 + i3] ^ bArr[i2]) & 255);
        i3++;
        i2++;
      }
      
      String y2 = String(bArr2, "UTF-8");
      return y2;
}
```

A possible transform function for the `script.js` 

```
function Transform() {
    var bArr = JSON.parse(base64.decode(Message))

    // JS implementation
    var length = bArr.length / 2;
    var length2 = bArr.length - length;
    var i2 = length2 - length;

    if (length % 2 != 0) {
        length = length - (length % 2)
    }

    var bArr2 = new Array(length);
    var i3 = 0;

    while (i3 < length) {
        bArr2[i3] = String.fromCharCode((hexStringToByteArray(bArr[length2 + i3]) ^ hexStringToByteArray(bArr[i2])) & 255);
        i3++;
        i2++;
    }

    return bArr2.join("")
}

function hexStringToByteArray(hexText) {
    var hexString = hexText.split("0x")[1]

    if (hexString.length == 1) { hexString = "0" + hexString }

    if (hexString.length % 2 !== 0) { throw "Must have an even number of hex digits to convert to bytes " + hexString }

    var numBytes = hexString.length / 2;
    var byteArray = new Uint8Array(numBytes);
    for (var i = 0; i < numBytes; i++) {

        if (hexText.includes("-")) {
            byteArray = parseInt(hexString.substr(i * 2, 2), 16) * -1;
        } else {
            byteArray[i] = parseInt(hexString.substr(i * 2, 2), 16);
        }

    }
    return byteArray;
}
```

## TODO

- [ ] Argument to change the replace method
- [ ] remove newlines
- [ ] Allow direct ingestion of APK files
- [ ] More efficient search process when running the execute command
- [ ] Auto frida script generation based on search rules
