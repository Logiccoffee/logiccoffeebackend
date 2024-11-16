package config

var ProfPost string = `
<div class="separator" style="clear: both; text-align: center;">
<a href="##PROFPIC##" style="margin-left: 1em; margin-right: 1em;">
<img border="0" height="400" src="##PROFPIC##" /></a>
</div><p></p><p></p>
<p>##BIO##</p>
`

var KatalogPost string = `
<div class="separator" style="clear: both; text-align: center;">
<a href="##URLCOVERBUKU##" style="margin-left: 1em; margin-right: 1em;">
<img border="0" height="400" src="##URLCOVERBUKU##" /></a>
</div><p></p><p></p><div class="alert-message alert">
<i class="fa fa-info-circle">
</i>Informasi Buku :<ul><li>ISBN : ##ISBN##</li>
<li>Terbit : ##TERBIT##</li>
<li>Ukuran : ##UKURAN##</li>
<li>Jumlah Halaman : ##JUMLAHHALAMAN##</li>
<li>Tebal : ##TEBAL##</li></ul></div>
<div class="alert-message success">
<i class="fa fa-file-text">
</i>Penulis :<ul>
##DAFTARPENULISDENGANTAGLI##
</ul>
<i class="fa fa-check-circle"></i>Editor :<ul><li>##EDITOR##</li></ul></div>
<blockquote class="tr_bq">##KALIMATPROMOSIBUKU##</blockquote>
<div><p><span class="firstcharacter">##HURUFPERTAMASINOPSIS##</span>
##SINOPSISBUKU##
</p></div><div>
<a class="button small buy" href="##LINKGRAMED##" style="color: white;">Gramedia</a>
<a class="button small download" href="##LINKPLAYBOOK##" style="color: white;">Google Play Books</a>
<a class="button small buy" href="##LINKKUBUKU##" style="color: white;">Kubuku</a>
<a class="button small download" href="##LINKMYEDISI##" style="color: white;">MyEdisi</a>
</div>
`

//##DAFTARPENULISDENGANTAGLI##
//<li>Muhammad Rizal Satria</li>
