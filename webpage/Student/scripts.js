document.addEventListener('DOMContentLoaded', function () {
    var timeModal = new bootstrap.Modal(document.getElementById('timeModal'));
    timeModal.show();
});

$(document).ready(function () {
    $(".headerpage").load("../components/header.html");


    let tempFiles = [];
    let uploadedFiles = [];

    $('#selectFileButton').click(function () {
        $('#fileInput').click();
    });

    $('#fileInput').change(function () {
        const files = Array.from(this.files);
        tempFiles = uploadedFiles.concat(files);
        updateFileList();
        $('#modal-file-list').toggle(tempFiles.length > 0);
    });

    $('#uploadModal').on('show.bs.modal', function () {
        tempFiles = [...uploadedFiles];
        updateFileList();
        $('#modal-file-list').toggle(tempFiles.length > 0);
    });

    $('#modal-file-list').on('click', '.delete-file', function () {
        const fileIndex = $(this).parent().data('index');
        tempFiles.splice(fileIndex, 1);
        updateFileList();
        $('#modal-file-list').toggle(tempFiles.length > 0);
    });

    $('#uploadConfirm').click(function () {
        uploadedFiles = [...tempFiles];
        updateMainFileList();
        $('#attachment-section').toggle(uploadedFiles.length > 0);
        tempFiles = [];
        updateFileList();
    });

    $('#file-list').on('click', '.delete-file', function () {
        const fileIndex = $(this).parent().data('index');
        uploadedFiles.splice(fileIndex, 1);
        updateMainFileList();
        $('#attachment-section').toggle(uploadedFiles.length > 0);
    });

    function updateMainFileList() {
        $('#file-list').empty();
        uploadedFiles.forEach((file, index) => {
            $('#file-list').append(`<div class="d-flex align-items-center mb-2 file-item" data-index="${index}">
                <span>${file.name}</span>
                <button class="btn btn-danger btn-sm ms-auto delete-file">刪除</button>
            </div>`);
        });
    }

    function updateFileList() {
        $('#modal-file-list').empty();
        tempFiles.forEach((file, index) => {
            $('#modal-file-list').append(`<div class="d-flex align-items-center mb-2 file-item" data-index="${index}">
                <span>${file.name}</span>
                <button class="btn btn-danger btn-sm ms-auto delete-file">刪除</button>
            </div>`);
        });
    }
});
