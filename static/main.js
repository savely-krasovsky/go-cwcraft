$('.search').focus();

new List('goods', {
    valueNames: ['name', 'id', 'type']
});

const clipboard = new ClipboardJS('.btn');

clipboard.on('success', function () {
    const copied = $('.copied');

    copied.css("visibility", "visible");

    copied.show();
    copied.fadeOut(1000);
});