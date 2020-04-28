# SMESH App keystore to CLIWallet json file convertor

SMESH keystore files are encrypted and cannot be used by CLIWallet.

CLIWallet wallet files are simple JSON files containing only the private and public keys.

This utility extracts all the private/public keypairs from a keystore and writes them to individual JSON files.

## Usage

smesh2cli -input <path to keystore> -password <SMESH password> -output <output file base>

* `input` e.g. `~/.config/spacemesh/my_wallet_0_2020-04-25T19-40-50.942Z.json`
* `password` - your smesh app password
* `output file base` e.g. `smdata/wallet`

A keystore file can contains multiple names key pairs. Assuming that you had two, workKey and funKey, they would be saved as two files `smdata/wallet0.json` and `smdata/wallet1.json`. The conversion process will indicate how they were saved.

## References

### Richard Moore's aes-js as used in the SMESH app

https://github.com/ricmoo/aes-js/blob/master/index.js#L613

This explained how to implement the line

`const aes1Ctr = new aes.ModeOfOperation.ctr(key, new aes.Counter(5)); // eslint-disable-line new-cap`

in https://github.com/spacemeshos/smapp/blob/853acca832aaf93b5d607c615dc44ab086dea3cd/app/infra/fileEncryptionService/fileEncryptionService.js#L57



### Go documents

https://godoc.org/crypto/cipher#example-NewCTR