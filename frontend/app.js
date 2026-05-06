const BASE_URL = 'http://localhost:4000';

async function apiCall(endpoint, panelId) {
    const panel = document.getElementById(panelId);
    
    // Toggle accordion logic
    const isAlreadyOpen = panel.classList.contains('open');
    
    // Close all panels first
    document.querySelectorAll('.panel').forEach(p => p.classList.remove('open'));
    
    if (!isAlreadyOpen) {
        panel.classList.add('open');
        panel.innerHTML = "<em>Sensing chakra... (Loading)</em>";

        try {
            const response = await fetch(`${BASE_URL}${endpoint}`);
            const data = await response.json();
            renderData(panel, data);
        } catch (error) {
            panel.innerHTML = "<span style='color:red;'>Failed to connect to the backend.</span>";
        }
    }
}

function renderData(container, data) {
    container.innerHTML = ""; 

    const mainDiv = document.createElement('div');
    mainDiv.className = 'holiday-card';

    if (data.message) {
        mainDiv.innerHTML = `<div class="message-box">${data.message}</div>`;
        if (data.is_holiday && data.occasion) {
            mainDiv.innerHTML += `<hr style="border: 0; border-top: 1px solid #9370db; margin: 10px 0;">
                                  <strong>${data.occasion}</strong>`;
        }
        container.appendChild(mainDiv);
        return;
    }

    if (data.holidays) {
        if (!data.holidays.length) {
            mainDiv.innerHTML = `<div class="message-box">No holidays next month</div>`;
            container.appendChild(mainDiv);
            return;
        }

        let content = "";
        data.holidays.forEach((h, index) => {
            content += `<strong>${h.occasion}</strong><br><small>${h.day}, ${h.date} ${h.year}</small>`;
            if (index < data.holidays.length - 1) {
                content += `<hr style="border: 0; border-top: 1px solid #9370db; margin: 10px 0; opacity: 0.3;">`;
            }
        });
        mainDiv.innerHTML = content;
        container.appendChild(mainDiv);
        return;
    }

    const list = data.occasions || data.dates || data.days;
    if (list) {
        mainDiv.innerHTML = `<ul style="margin: 0; padding-left: 20px;">${list.map(item => `<li>${item}</li>`).join('')}</ul>`;
        container.appendChild(mainDiv);
        return;
    }

    if (data.found === false || (data.found && data.holiday)) {
        if (data.found === false) {
            mainDiv.innerHTML = `<div class="message-box">${data.message}</div>`;
        } else {
            mainDiv.innerHTML = `
                <strong>${data.holiday.occasion}</strong>
                <p>Coming up in <strong>${data.days_away} days</strong>!</p>
                <small>${data.holiday.date}, ${data.holiday.year}</small>`;
        }
        container.appendChild(mainDiv);
    }
}