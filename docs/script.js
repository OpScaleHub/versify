document.addEventListener('DOMContentLoaded', () => {
    // Smooth scrolling for nav links
    const navLinks = document.querySelectorAll('.nav-links a');
    for (const link of navLinks) {
        if (link.getAttribute('href').startsWith('#')) {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const targetId = link.getAttribute('href');
                const targetElement = document.querySelector(targetId);

                if (targetElement) {
                    window.scrollTo({
                        top: targetElement.offsetTop - 70, // Adjusted for fixed header
                        behavior: 'smooth'
                    });
                }
            });
        }
    }

    // Installation tabs
    const tabLinks = document.querySelectorAll('.tab-link');
    const tabContents = document.querySelectorAll('.tab-content');

    tabLinks.forEach(link => {
        link.addEventListener('click', () => {
            const os = link.dataset.os;

            tabLinks.forEach(l => l.classList.remove('active'));
            link.classList.add('active');

            tabContents.forEach(c => c.classList.remove('active'));
            document.getElementById(`${os}-instructions`).classList.add('active');
        });
    });

    // OS detection for default tab
    const platform = navigator.platform.toLowerCase();
    let defaultOS = 'linux';
    if (platform.includes('mac') || platform.includes('darwin')) {
        defaultOS = 'macos';
    } else if (platform.includes('win')) {
        defaultOS = 'windows';
    }

    const defaultTab = document.querySelector(`.tab-link[data-os=${defaultOS}]`);
    if (defaultTab) {
        defaultTab.click();
    }


    // Populate installation commands
    const commands = {
        'linux-amd64': `curl -L "https://github.com/OpScaleHub/versify/releases/latest/download/versify-linux-amd64.tar.gz" -o versify.tar.gz
tar -xzvf versify.tar.gz
chmod +x versify-linux-amd64
sudo mv versify-linux-amd64 /usr/local/bin/versify
rm versify.tar.gz`,
        'linux-arm64': `curl -L "https://github.com/OpScaleHub/versify/releases/latest/download/versify-linux-arm64.tar.gz" -o versify.tar.gz
tar -xzvf versify.tar.gz
chmod +x versify-linux-arm64
sudo mv versify-linux-arm64 /usr/local/bin/versify
rm versify.tar.gz`,
        'macos-amd64': `curl -L "https://github.com/OpScaleHub/versify/releases/latest/download/versify-darwin-amd64.tar.gz" -o versify.tar.gz
tar -xzvf versify.tar.gz
chmod +x versify-darwin-amd64
sudo mv versify-darwin-amd64 /usr/local/bin/versify
rm versify.tar.gz`,
        'windows-amd64': `$url = "https://github.com/OpScaleHub/versify/releases/latest/download/versify-windows-amd64.zip"
$output = "versify.zip"
Invoke-WebRequest -Uri $url -OutFile $output
Expand-Archive -Path $output -DestinationPath .
# Move to a directory in your PATH (e.g., C:\\Windows\\System32 or a custom one)
# Make sure to run PowerShell as Administrator for this step
Move-Item -Path ".\\versify-windows-amd64.exe" -Destination "C:\\Windows\\System32\\versify.exe"
Remove-Item -Path $output`
    };

    for (const id in commands) {
        const el = document.getElementById(`${id}-code`);
        if (el) {
            el.textContent = commands[id];
        }
    }
});