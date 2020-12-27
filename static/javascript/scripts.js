(function() {
    var burger = document.querySelector('.burger');
    var menu = document.querySelector('#' + burger.dataset.target);
    burger.addEventListener('click', function() {
        burger.classList.toggle('is-active');
        menu.classList.toggle('is-active');
    });

    function getAll(selector) {
        var parent = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : document;

        return Array.prototype.slice.call(parent.querySelectorAll(selector), 0);
    }

    function closeDropdowns() {
        $dropdowns.forEach(function($el) {
            $el.classList.remove('is-active');
        });
    }
    var rootEl = document.documentElement;
    var $modals = getAll('.modal');
    var $modalButtons = getAll('.modal-button');
    var $modalCloses = getAll('.modal-background, .modal-close, .modal-card-head .delete, .modal-card-foot .button');

    if ($modalButtons.length > 0) {
        $modalButtons.forEach(function($el) {
            $el.addEventListener('click', function() {
                var target = $el.dataset.target;
                openModal(target);
            });
        });
    }

    if ($modalCloses.length > 0) {
        $modalCloses.forEach(function($el) {
            $el.addEventListener('click', function() {
                closeModals();
            });
        });
    }

    function openModal(target) {
        var $target = document.getElementById(target);
        rootEl.classList.add('is-clipped');
        $target.classList.add('is-active');
    }

    function closeModals() {
        rootEl.classList.remove('is-clipped');
        $modals.forEach(function($el) {
            $el.classList.remove('is-active');
        });
    }

    document.addEventListener('keydown', function(event) {
        var e = event || window.event;
        if (e.keyCode === 27) {
            closeModals();
            closeDropdowns();
        }
    });
})();

var showErrorMessage = function(message) {
    var errorEl = document.getElementById("error-message")
    errorEl.textContent = message;
    errorEl.style.display = "block";
};

var handleFetchResult = function(result) {
    if (!result.ok) {
        return result.json().then(function(json) {
            if (json.error && json.error.message) {
                throw new Error(result.url + ' ' + result.status + ' ' + json.error.message);
            }
        }).catch(function(err) {
            showErrorMessage(err);
            throw err;
        });
    }
    return result.json();
};

var handleResult = function(result) {
    if (result.error) {
        showErrorMessage(result.error.message);
    }
};

var createCheckoutSession = function(priceID) {
    return fetch("/billing/create-checkout-session", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            priceID: priceID
        })
    }).then(function(result) {
        return result.json();
    });
}

fetch("/billing/setup")
    .then(handleFetchResult)
    .then(function(json) {
        var publishableKey = json.publishableKey;
        var smallPrice = json.smallPrice;
        var mediumPrice = json.mediumPrice;
        var largePrice = json.largePrice;

        var stripe = Stripe(publishableKey);
        document.getElementById("smallPlan")
            .addEventListener("click", function(event) {
                createCheckoutSession(smallPrice).then(function(data) {
                    stripe
                        .redirectToCheckout({
                            sessionId: data.sessionID
                        })
                        .then(handleResult);
                })
            })

        document.getElementById("mediumPlan")
            .addEventListener("click", function(event) {
                createCheckoutSession(mediumPrice).then(function(data) {
                    stripe
                        .redirectToCheckout({
                            sessionId: data.sessionID
                        })
                        .then(handleResult);
                })
            })

        document.getElementById("largePlan")
            .addEventListener("click", function(event) {
                createCheckoutSession(largePrice).then(function(data) {
                    stripe
                        .redirectToCheckout({
                            sessionId: data.sessionID
                        })
                        .then(handleResult);
                })
            })

    })