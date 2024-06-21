const productSection = document.getElementById('bouquets');
const modal = document.getElementById('modal');
const modalTitle = document.getElementById('modal-title');
const modalDescription = document.getElementById('modal-description');
const orderForm = document.getElementById('order-form');

orderForm.addEventListener('submit', async function(event) {
    event.preventDefault();
    console.log('Форма заказа отправлена.');

    const name = document.getElementById('name').value;
    const phone = document.getElementById('phone').value;
    const bouquetName = modalTitle.textContent;
    console.log('Данные для отправки на сервер:', {
        bouquetName: bouquetName,
        customerName: name,
        customerPhone: phone,
    });

    console.log('Отправка запроса на сервер для оформления заказа...');

    try {
        const response = await fetch('http://localhost:8080/api/orders', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                bouquetName: bouquetName,
                customerName: name,
                customerPhone: phone,
            }),
        });

        if (!response.ok) {
            throw new Error('Отклик сети был не в порядке');
        }

        console.log('Перед ожиданием ответа сервера.');
        const data = await response.text();
        console.log('Ответ сервера:', data);
        event.preventDefault();
        console.log('Заказ успешно оформлен.');
        modal.style.display = 'none';
        window.location.href = "cart.html";
    } catch (error) {
        console.error('Ошибка при оформлении заказа:', error);
        console.log('Ошибка в fetch:', error);
    }
});

async function initialize() {
    try {
        const response = await fetch('http://localhost:8080/api/products');
        const products = await response.json();
        displayProducts(products);
    } catch (error) {
        console.error('Ошибка загрузки товаров:', error);
    }
}

function displayProducts(products) {
    products.forEach(product => {
        const card = document.createElement('div');
        card.className = 'product-card';

        card.innerHTML = `
            <img src="${product.image}" alt="">
            <div class="product-name">
                <div class="text1">${product.name}</div>
            </div>
            <div class="text3">${product.description}</div>
            <button class="product-btn" data-original-text="${product.price}">${product.price}</button>
        `;

        productSection.appendChild(card);

        const button = card.querySelector('.product-btn');
        button.addEventListener('mouseover', function() {
            button.textContent = 'Заказать';
        });

        button.addEventListener('mouseout', function() {
            button.textContent = button.getAttribute('data-original-text');
        });

        const closeButton = document.getElementsByClassName('close')[0];
        closeButton.addEventListener('click', function() {
            modal.style.display = 'none';
        });

        window.addEventListener('click', function(event) {
            if (event.target == modal) {
                modal.style.display = 'none';
            }
        });

        button.addEventListener('click', function() {
            modalTitle.textContent = product.name;
            modalDescription.textContent = product.description;
            modal.style.display = 'block';
        });
    });
}

initialize();
