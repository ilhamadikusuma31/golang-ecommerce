<h1>Project Golang Web e-commerce </h1>
<hr>
<h3>tech stack</h3>
<ul>
    <li>Golang</li>
    <li>mongodb</li>
    <li></li>
</ul>
<hr>
<h3>Penjelasan</h3>
<ol>
    <li>
        <h4><b>Aggregate</b></h4>
        <p>operasi di mongodb misal where select avg sum di sql</p>
        <p>misalkan disini ada database bernama <b>belajar-mongodb-copass</b> dan ada collection/table person</p>
        <img src="info/collection.png" alt="tes">
        <hr>
        <ul>
            <li>
                <h4>match</h4>
            </li>
            <p>kalo match itu adalah mereturnkan semua document/row yang cocok</p>
            <p>dalam studi kasus ini adalah yang isActive nya false</p>
            <img src="info/syntax-where.png" alt="">
            <img src="info/collection-match.png" alt="">
            <li>
                <h4>group</h4>
            </li>
            <p>kalo misalnya tadi sudah <i>di-match</i> maka group akan mengembalikan dokumen yang serupa seperti match
                hanya saja fieldnya sudah ditentukan</p>
            <p>dalam studi kasus ini adalah yang isActive di-match-kan lalu di-group berdasarkan name</p>
            <p>untuk sintaks dari group yaitu field yang ingin diambil maka akan menjadi key <i>_id</i></p>
            <img src="info/collection-match-group.png" alt="">
            <li>
                <h4>unwind</h4>
            </li>
            <p>membuat dokumen/row baru dari setiap dokumen/row yang mempunyai nested value</p>
            <p>studi kasus ini setiap dokumen punya field tags nah lalu nanti setiap tags akan dibuatkan dokumen/row
                masing-masing</p>
            <img src="info/collection-unwind.png" alt="">
        </ul>
    </li>
</ol>