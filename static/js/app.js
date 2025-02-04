function Prompt() {
    let toast = function(c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c;
        const Toast = Swal.mixin({
        toast: true,
        title: msg,
        position: position,
        icon: icon,
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
            toast.onmouseenter = Swal.stopTimer;
            toast.onmouseleave = Swal.resumeTimer;
        }
        });
        Toast.fire({});
    }

    let success = function(c) {
        const {
            icon = "success",
            msg = "",
            title = "",
            href = "index.html",
            linktxt = "Link!"
        } = c;
        Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: "<a href=" + href + ">" + linktxt + "</a>"
            });
    }

    async function custom(c){
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,

        } = c;

        const { value: result } = await Swal.fire({
            icon : icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined){
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        })
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !=="") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callbacl(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        custom: custom,
    }
}

function CheckAvail(c) {
    // const {room_id = "0"} =c;
    
    let html =`
	<form id="check-availability-form" action="" method="post" novalidate class="needs-validation"> 
		<div class="row" id="reservation-dates-modal">
			<div class="col-md-6">
				<input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
			</div>
			<div class="col-md-6">
				<input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
			</div>
		</div>
	</form>
	`
	attention.custom({
		msg: html, 
		title: "Search Availability",
		
		willOpen: () => {
			const elem = document.getElementById('reservation-dates-modal');
			const rp = new DateRangePicker(elem, {
				format: 'yyyy-mm-dd',
				minDate: new Date(),
				showOnFocus: true,
				
			})
            //console.log(">>>>>>>>>>>willOpen Ran..");
            //console.log("console log test in JS");
            //console.log("CRSF in js >", csrf_t);
            
            
		},

		didOpen: () => {
			document.getElementById("start").removeAttribute('disabled');
			document.getElementById("end").removeAttribute('disabled');
            //console.log("didOpen Ran..")
		},

		callback: function(result) {
			
            let form = document.getElementById("check-availability-form");
			let formData = new FormData(form);
			formData.append("csrf_token", csrf_t);
			formData.append("room_id", roomID);
            //console.log("callback Ran..");

			fetch('/search-availability-json', {
				method: "post",
				body: formData,
			}) 
                
                .then(response => response.json())
				.then(data => {
					if (data.ok) {
						attention.custom({
							icon: 'success',
							showConfirmButton: false,
							msg: '<p>Room is Available!</p>'
								+ '<p><a href="/book-room?id='
								+ data.room_id
								+ '&s='
								+ data.start_date
								+ '&e='
								+ data.end_date
								+ '" class="btn btn-primary">'
								+ 'Book now!</a></p>',
						})

					} else {
						attention.custom({
							icon: 'error',
							msg: "No availibility",
							showConfirmButton: false,
						})
					}
				})
		}
	});
}