# Goroutines with Web scraping

## Demo

This demonstration consist of a dictionary file that contain all the english words (about 300 k). For each word will be made a request to the linguee website.

An example url can be *https://www.linguee.es/ingles-espanol/search?query=potato*

From the request are extracted 20 phrases that contain it in context, and its translation into spanish.

Phrases obtained will be saved in an xml file, when the scraper finished.

```Xml file
<dictionary>
    <words>
        <word>
            <value>anoisgabnoiwga</value>
            <examples></examples>
        </word>
        <word>
            <value>hello</value>
            <example>
                <english>I therefore believe that a standardised portal along amazon lines -  partly, but not exclusively, computer-driven and greeting  people with the words, Hello, you are now at our Brussels premises. </english>
                <spanish>Por ello, creo que un portal normalizado similar a amazon -que funcione en parte, pero no  exclusivamente, por ordenador y que salude a las  personas con las palabras «Hola, te encuentras en nuestras oficinas de Bruselas. </spanish>
            </example>
            <example>
                <english>Hello, Just a word to let you know how satisfied I am with my internet banking account. </english>
                <spanish>Hola: Sólo unas palabras para hacerles saber lo satisfecho que estoy con mi cuenta en la banca por internet. </spanish>
            </example>
        </word>
    </words>
</dictionary>
```

Web requests are controlled by a semaphore channel that limits the number of concurrent requests, its buffer is the value of your cpu number, *runtime.NumCPU()*

## Limits

Linguee website do an ip ban if it receives many petitions, from this the following requests receive a 403 error.

## Libraries

External used libraries are:

- goquery : Search the phrases into the html
- iconv.v1: Linguee is in ISO-8859-15, this library converts it into UTF-8. NOTE: For compile it in windows you can install mingw
