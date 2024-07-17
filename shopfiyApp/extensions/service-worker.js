self.addEventListener('push', function (event) {
    const data = event.data.json();
    const options = {
      body: data.body,
    //   icon: 'icon.png', // Optional: path to an icon
    //   badge: 'badge.png' // Optional: path to a badge
    };
    event.waitUntil(
      self.registration.showNotification(data.title, options)
    );
  });
  