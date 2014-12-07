The project is currently hosted over at [GitHub](https://github.com/lesocr). It is not public, in order to help other people to innovate. We might make it public once the final milestone is reached by our team.

However, you can upload a file here and we will do our best to parse it, so you can have a peek at the capabilities of our OCR. We accept JPEG, PNG and BMP files, with a file size under 512 KiB.

<form method="POST" enctype="multipart/form-data">
	<input type="file" name="file" />
	<a href="#" onclick="document.forms[0].submit()">Submit</a>
</form>
<br />

{{.OCR_RESULT}}
