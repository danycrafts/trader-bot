from playwright.sync_api import sync_playwright

def verify_check_balance():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        # Use a larger viewport to ensure the button is visible
        page = browser.new_page(viewport={'width': 1280, 'height': 1024})
        try:
            # Vite default port is usually 5173
            page.goto("http://localhost:5173")

            # Wait for the app to load
            page.wait_for_selector("#App")

            # Check for the new button "Check Balance"
            # It should be visible
            button = page.locator("button", has_text="Check Balance")
            if button.is_visible():
                print("Button 'Check Balance' is visible.")
            else:
                print("Button 'Check Balance' is NOT visible.")
                exit(1)

            # Take a screenshot
            page.screenshot(path="/home/jules/verification/verification.png")
            print("Screenshot taken.")

        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()

if __name__ == "__main__":
    verify_check_balance()
