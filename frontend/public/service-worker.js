// public/service-worker.js

self.addEventListener('push', function(event) {
    // convert string to JSON
    const data = event.data.json();
    const title = data.title;
    
    const options = {
        body: data.body,
        icon: data.icon,
        data : {
            url : data.url
        },
    };
    event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener('notificationclick', function(event) {
    clients.openWindow(event.notification.data.url);
}, false);

