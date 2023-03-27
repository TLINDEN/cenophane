package cmd

const formtemplate = `
<!DOCTYPE html>
<!-- -*-web-*- -->
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="description" content="upload form" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>File upload form</title>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js"></script>

  </head>
  <body>
    <h4>Upload form {{ .Id }}</h4>
    <!-- Status message -->
    <div class="statusMsg"></div>

    <!-- File upload form -->
    <div class="col-lg-12">
      <form id="fupForm" enctype="multipart/form-data" action="/v1/uploads" method="POST">
        <div class="form-group">
          <label for="expire">Expire</label>
          <input type="expire" class="form-control" id="expire" name="expire" placeholder="Enter expire"/>
        </div>
        <div class="form-group">
          <label for="file">Files</label>
          <input type="file" class="form-control" id="file" name="uploads[]" multiple />
        </div>
        <input type="submit" name="submit" class="btn btn-success submitBtn" value="Upload"/>
      </form>
    </div>


    <script>
     $(document).ready(function(){
       // Submit form data via Ajax
       $("#fupForm").on('submit', function(e){
         e.preventDefault();
         $.ajax({
           type: 'POST',
           url: '/v1/uploads',
           data: new FormData(this),
           dataType: 'json',
           contentType: false,
           cache: false,
           processData:false,
           beforeSend: function(xhr){
               $('.submitBtn').attr("disabled","disabled");
               $('#fupForm').css("opacity",".5");
               xhr.setRequestHeader('Authorization', 'Bearer {{.Id}}');
           },
           success: function(response){
             $('.statusMsg').html('');
             if(response.success){
                 $('#fupForm')[0].reset();
                 $('.statusMsg').html('<p class="alert alert-success">Your upload is available at <a href="'
                                      +response.uploads[0].url+'">here</a> for download</p>');
                 $('#fupForm').hide();
             }else{
               $('.statusMsg').html('<p class="alert alert-danger">'+response.message+'</p>');
             }
             $('#fupForm').css("opacity","");
             $(".submitBtn").removeAttr("disabled");
           }
         });
       });
     });
    </script>
  </body>
</html>

    
`
