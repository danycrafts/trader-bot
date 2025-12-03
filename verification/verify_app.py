from playwright.sync_api import sync_playwright

def verify_frontend():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        page.on("console", lambda msg: print(f"Browser Console: {msg.text}"))
        page.on("pageerror", lambda exc: print(f"Browser Error: {exc}"))

        # Mock window.go object which Wails usually injects
        page.add_init_script("""
            window.go = {
                main: {
                    App: {
                        Login: (email, password) => {
                            console.log("Login called");
                            return Promise.resolve({id: 1, email: email});
                        },
                        Register: (email, password) => {
                             console.log("Register called");
                             return Promise.resolve();
                        },
                        GetAccountBalance: () => Promise.resolve(10000.50),
                        GetPortfolio: () => Promise.resolve([
                            {Symbol: "AAPL", Qty: 10, EntryPrice: 150.0, CurrentPrice: 155.0, PnL: 50.0},
                            {Symbol: "TSLA", Qty: 5, EntryPrice: 700.0, CurrentPrice: 690.0, PnL: -50.0}
                        ]),
                        GetRecentTrades: () => Promise.resolve([
                             {id: "1", Symbol: "AAPL", Side: "buy", Strategy: "Momentum", EntryPrice: 150.0, PnL: 10.0},
                             {id: "2", Symbol: "TSLA", Side: "sell", Strategy: "MeanReversion", EntryPrice: 700.0, PnL: -5.0}
                        ]),
                        SearchStocks: (query) => Promise.resolve(["AAPL", "GOOGL", "MSFT"]),
                        PlaceOrder: (s, q, side) => Promise.resolve(),
                        StartMarketStream: (s) => Promise.resolve("Stream Started"),
                        GetSettings: (id) => Promise.resolve({
                            user_id: 1,
                            alpaca_api_key: "test",
                            theme: "light",
                            notifications_email: true
                        }),
                        SaveSettings: (s) => Promise.resolve()
                    }
                }
            };
        """)

        try:
            # Navigate to the app
            page.goto("http://localhost:5173")
            page.wait_for_load_state("networkidle")

            # Screenshot Login
            page.screenshot(path="verification/login_page.png")
            print("Login page screenshot taken.")

            # Perform Login
            page.fill('input[name="email"]', "test@example.com")
            page.fill('input[name="password"]', "password")
            page.click('button:has-text("Sign In")')

            # Wait for navigation
            page.wait_for_timeout(1000)

            # Screenshot Dashboard
            page.screenshot(path="verification/dashboard_page.png")
            print("Dashboard page screenshot taken.")

            # Navigate to Market
            page.click('div[role="button"]:has-text("Market")')
            page.wait_for_timeout(1000)
            page.screenshot(path="verification/market_page.png")
            print("Market page screenshot taken.")

            # Navigate to Settings
            page.click('div[role="button"]:has-text("Settings")')
            page.wait_for_timeout(1000)
            page.screenshot(path="verification/settings_page.png")
            print("Settings page screenshot taken.")

        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()

if __name__ == "__main__":
    verify_frontend()
