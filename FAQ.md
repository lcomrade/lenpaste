# Frequently Asked Questions
## Why such a name?
In the USSR, some organizations were named according to the principle:
abbreviated name of city + main activity of the organization.
For example: LenFilm, MosGas.
Lenpaste repeats this tradition.


## What are the system requirements?
There aren't any. In fact, Lenpaste is so lightweight that you can run it almost anywhere.
At idle time the container uses about 10 MB of RAM, it uses almost no CPU,
and the disk load is also minimal (all HTML pages are loaded into memory at startup).


## Why isn't the source code posted on GitHub?
There are several reasons why I don't use GitHub,
all of which are detailed [here](https://git.lcomrade.su/root/give-up-github).
Here I will only describe the basic idea.

GitHub has gained too much power over open source developers.
It's time to take that power away from the capitalists.


## Does Lenpaste collect telemetry?
**No.** No data about users or the server is collected.
We respect your privacy!


## I want to modify the files in the `./web/` dir, do I have to provide the source code?
No, they shouldn't, because these files are distributed without a license.
That is, the AGPLv3 license does not apply to them.

But remember that you may violate AGPLv3 item 13 about explicitly providing source code for the program.
**You can't remove the link to download the source code from Lenpaste!**
If Lenpaste doesn't offer to download its source code, it's a license violation.


## Why is there no support for encryption of paste?
Because if the server is compromised, malicious JavaScript can be downloaded to the user's devices, which will compromise the encryption.
If you need encryption, then use it manually with GPG or something similar.


##  Why can't the paste be password-protected?
For the same reason that encryption is not supported in the browser.
It's not secure.


## Who develops Lenpaste?
Basically, I develop a project alone.
Also some people help with translation, testing, translation and so on.
Full list of contributors here: https://paste.lcomrade.su/about/authors


## Do you need any help developing Lenpaste?
**Yes.** You can do the following:
- Create logo (favicon) to Lenpaste.
- Translation into other languages.
- Run your server and add it to [Lenmonitor](https://monitor.lcomrade.su/add).
- Write a tutorial, a review, or make a video.
- Suggest an idea.
- Report a bug.
- And so on.

If you still have questions, please contact me: Leonid Maslakov <root@lcomrade.su>.
