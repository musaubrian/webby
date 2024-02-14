# webby

Easily test your typeform webhook connection

**Why** - TypeForm doesn't let you setup a webhook to your form and connect it to a loccally running app.

Deploy this on a server somewhere, point the webhook to that and presto, you can view the data your form sends.

> You need to have go installed

```sh
git clone https://github.com/musaubrian/webby

cd webby

go build .

./webby

```
