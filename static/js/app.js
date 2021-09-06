// Prompt is a js module for all alerts, notifications and custom pop-ups dialogs
function Prompt() {
    let toast = function (c) {

        const {
            msg = "",
            icon = "success",
            position = "top-end",
            timer = 3000
        } = c;

        // ref: https://sweetalert2.github.io/ and search for Toast
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            icon: icon,
            position: position,
            showConfirmButton: false,
            timer: timer,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let icon = function (c) {

        const {
            msg = "",
            icon = "success",
            title = "",
            footer = "",
        } = c;

        const Success = Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: footer
        })
    }

    // pop-up for error msg
    let error = function (c) {

        const {
            msg = "",
            icon = "error",
        } = c;

        const Err = Swal.fire({
            icon: icon,
            title: 'Sorry',
            text: msg,
          })
    }

    // for multiple inputs using sweetalert2
    // await can only be used with async func
    async function custom(c) {

        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton=true,
        } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,

            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
            // open the calender when the user clicked on the textbox
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        })

        // this is for booking rooms and checking availability of the rooms in the background
        if (result) {
            // result.dismiss is not exactly(!==) to Swal.DismissReason.canel i.e
            // if they didn't hit the cancel button then
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    // calling a callback on client page when someone fills a form
                    if (c.callback !== undefined) {
                        c.callback(result);
                    } 
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        icon: icon,
        error: error,
        custom: custom,
    }
}