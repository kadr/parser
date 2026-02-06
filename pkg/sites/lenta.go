package sites

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	urlutils "net/url"
	"os"
	"strings"
)

type price struct {
	Cost  int `json:"cost"`
	Price int `json:"price"`
}

type item struct {
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Prices price  `json:"prices"`
	Slug   string `json:"slug"`
}

type SiteResponse struct {
	Items []item `json:"items"`
}

type Categories struct {
	Categories []Category `json:"categories"`
}
type Category struct {
	Id   int
	Name string
	Slug string
}

type Lenta struct{}

const (
	HOST string = "https://lenta.com/"
)

func NewLenta() *Lenta {
	return &Lenta{}
}

func (l *Lenta) Parse(category string, limit *int) ([]Result, error) {
	newLimit := 100
	if limit != nil {
		newLimit = *limit
	}
	url := fmt.Sprintf("%s/api-gateway/v1/catalog/items", HOST)
	categoryId, err := getCategoryID(category)
	if err != nil {
		fmt.Println(err)
		return []Result{}, err
	}
	filter := map[string]any{"categoryId": categoryId, "limit": newLimit, "offset": 0}
	var result SiteResponse
	body, err := makeRequest(http.MethodPost, url, &filter)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return []Result{}, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&result)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return []Result{}, err
	}

	return mapping(result, category), nil
}

func getCategoryID(category string) (int, error) {
	url := fmt.Sprintf("%s/api-gateway/v1/catalog/categories", HOST)
	resp, err := makeRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return 0, err
	}
	defer resp.Close()
	var categories Categories
	err = json.NewDecoder(resp).Decode(&categories)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return 0, err
	}
	for _, cat := range categories.Categories {
		if cat.Name == category || strings.Contains(cat.Name, category) {
			return cat.Id, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Category: %s not found", category))
}

func makeRequest(method, url string, filter *map[string]any) (io.ReadCloser, error) {
	var err error
	client := &http.Client{}
	if proxy, exist := os.LookupEnv("PROXY"); exist {
		proxy = strings.Trim(proxy, " ")
		if len(proxy) > 0 {
			parsedProxy, err := urlutils.Parse(proxy)
			if err == nil {
				transport := &http.Transport{Proxy: http.ProxyURL(parsedProxy)}
				client = &http.Client{Transport: transport}
			} else {
				fmt.Println("Не удалось распарсить проху")
			}
		}
	}

	var req *http.Request
	if filter != nil {
		newFilter, err := json.Marshal(filter)
		if err != nil {
			fmt.Println(err.Error())
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(newFilter))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return nil, err
	}
	for key, val := range getHeaders() {
		req.Header.Add(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		res, _ := io.ReadAll(resp.Body)
		return nil, errors.New(fmt.Sprintf("Сервер вернул ошибку. Статус: %d. %s", resp.StatusCode, string(res)))
	}
	return resp.Body, nil
}

func mapping(response SiteResponse, category string) []Result {
	var result []Result
	for _, resp := range response.Items {
		linkToProduct := fmt.Sprintf("%s/product/%s", HOST, resp.Slug)
		result = append(result, Result{Name: resp.Name, Category: category, Link: linkToProduct, Price: float64(resp.Prices.Price)})
	}

	return result
}

func getHeaders() map[string]string {
	cookie := os.Getenv("LENTA_COOKIE")
	cookie = strings.Trim(cookie, " ")
	return map[string]string{
		"Host":                "lenta.com",
		"User-Agent":          "Mozilla/5.0 (X11; Linux x86_64; rv:148.0) Gecko/20100101 Firefox/148.0",
		"Accept":              "application/json",
		"Accept-Language":     "ru,en-US;q=0.9,en;q=0.8",
		"Content-Type":        "application/json",
		"Content-Length":      "139",
		"Referer":             "https://lenta.com/catalog/avtotovary-18199/",
		"SessionToken":        "7FA41F71766A32D0D2E6596334F46D9A",
		"DeviceID":            "d3558ef8-445f-ca5e-de21-601065a9f5a2",
		"X-Retail-Brand":      "lo",
		"X-Platform":          "omniweb",
		"X-Device-ID":         "d3558ef8-445f-ca5e-de21-601065a9f5a2",
		"X-Device-OS":         "Web",
		"X-Device-OS-Version": "12.4.8",
		"X-Delivery-Mode":     "pickup",
		"X-Domain":            "moscow",
		"Client":              "angular_web_0.0.2",
		"X-User-Session-Id":   "51d22be6-44d1-98f6-c50f-6df30dec29a5",
		"Experiments":         "exp_recommendation_cms.true, exp_lentapay.test, exp_profile_bell.test, exp_cl_omni_authorization.test, exp_fullscreen.test, exp_onboarding_editing_order.test, exp_cart_new_carousel.default, exp_sbp_enabled.default, exp_profile_settings_email.default, exp_cl_omni_refusalprintreceipts.test, exp_search_suggestions_popular_sku.default, exp_cl_new_csi.default, exp_cl_new_csat.default, exp_delivery_price_info.test, exp_interval_jump.test, exp_cardOne_promo_type.test, exp_birthday_coupon_skus.test, exp_qr_cnc.test, exp_where_place_cnc.test, exp_editing_cnc_onboarding.test, exp_editing_cnc.test, exp_pickup_in_delivery.test, exp_welcome_onboarding.default, exp_where_place_new.default, exp_start_page.default, exp_default_payment_type.default, exp_start_page_onboarding.default, exp_search_new_logic.default, exp_referral_program_type.default, exp_new_action_pages.default, exp_items_by_rating.test, exp_can_accept_early.default, exp_online_subscription.default, exp_hide_cash_payment_for_cnc_wo_adult_items.default, exp_prices_per_quantum.test, exp_web_chips_online.test, exp_chips_online.default, exp_promo_without_benefit.default, exp_cart_forceFillDelivery.default, exp_banner_sbp_checkout_step_3.control, exp_kit_banner_sbp_checkout_step_3.default, exp_kit_badge_sbp_checkout_step_3.default, exp_profile_stories.test, exp_cl_new_ui_csi_comment.default, exp_in_app_update.default, exp_sorting_catalog.default, exp_aa_test_2025_04.test, exp_product_page_by_blocks.Default, exp_without_a_doorbell_new.default, exp_edit_payment_type.test, exp_edit_payment_type_new.test, exp_search_photo_positions.default, exp_new_matrix.test, exp_another_button_ch.default, exp_progressbar_and_title.test, exp_auto_fill_coupon.default, exp_promo_and_bonus.test, exp_about_cnc_optimization.default, exp_online_categories.default, exp_no_intervals.default, exp_web_b2b_excel_load.default, exp_cart_save_with_promo.default, exp_email_optional_full_registration.default, exp_cl_new_rateapp.default, exp_similar_goods_cart.test, exp_cart_redesign_promocode.default, exp_search_new_filters.default, exp_loyalty_categories_labels.default, exp_search_multicard.test, exp_delivery_promocode_bd_coupon.default, exp_search_disable_fuzziness.default, exp_ui_catalog_level_2.default, exp_fullscreen_inapp_vs_native.test1, exp_search_collections_ranking.default, exp_search_elastic_tokens.default, exp_cl_new_tapbar.default, exp_cl_new_tapbar_tab.default, exp_cart_free_sample.default, exp_personal_promo_detail_for_delivery.default, exp_search_combined_field.default, exp_search_unified.default, exp_web_personal_promo_detail_for_delivery.default, exp_web_personal_promo_delivery_chips.control, exp_b2b_web_mob_checkout.default, exp_personal_promo_delivery_chips.default, exp_ds_cnc_pers_recom.test, exp_ds_mntk_pers_recom.default, exp_shopping_statistics.default, exp_pin_create_button.default, exp_search_ui_catalog_pim.default, exp_search_video.default, exp_search_pinned_reviews.default, exp_sbp_instead_of_lenta_pay.default, exp_card1_start_page.default, exp_status_assemble_completed.test, exp_cl_new_ui_csi_comment2.default, exp_online_subscription_discount.default, exp_start_page_button_notifications.default, exp_quick_checkout.default, exp_quick_checkout_update.default, exp_search_no_stock.true, exp_main_page_new_mode_shop.default, exp_brief_description_promo.default, exp_new_offer_new_user_v1.default, exp_order_feedback_show.default, exp_leave_order_at_door.test, exp_leave_order_at_door_new.test, exp_search_quantity_discount_promo.default, exp_start_page_button_navigation_off.default, exp_obi_webview.true, exp_huawei_adjust_new_tokens.true, exp_import_goods_in_basket.efault, exp_unpin_tabbar.default, exp_mna_orders_editing.default, exp_consent.default, exp_main_stories.test, exp_from_store_myself.default, exp_new_bs_catalog_startpage.default, exp_be_soon_show_explain_message.default, exp_startpage_mainpage_new_address_design.default, exp_bubble_discount_startpage_mainpage.test, exp_startpage_zone_description.default, exp_ds_pd_carousel.default, exp_ds_pers_recom_delivery_2.test, exp_search_ds_catboost_2.control, exp_new_user_promo_profile.default, exp_novikov_test.OFF, exp_order_created_action_banner.default, exp_ds_mntk_pers_cat.default, exp_search_ds_empty_recom.default, exp_badges_pers_cashback.default, exp_temp_exp_ds_pd_carousel_android.default, exp_temp_exp_ds_pd_carousel_android_general.default, exp_temp_samesplit_check_f.default, exp_interval_jump_30.default, exp_temp_exp_ds_pd_carousel_ios_general.default, exp_search_purchased_badge.default, exp_pwa_cart.default, exp_pwa_checkout.default, exp_auto_detection_store_for_new_user.default, exp_return_available_items.default, exp_b2c_onboarding_send_cart.default, exp_b2b_send_cart.default, exp_b2c_send_cart.default, exp_cart_item_modify_version20.default, exp_auth_sber_id.default, exp_search_voice_search_ai.default, exp_startpage_redesign_qr_and_loy.default, exp_web_aa_2026_01_v1.test, exp_startpage_redesign_missions.default, exp_search_pdp_big_photo.default, exp_startpage_tab_shop_on_the_map.default, exp_startpage_logics_button_pickup.default, exp_open_screen_card1_profile_without_address.default, exp_web_cancel_to_edit_cnc.default, exp_ch_how_much_unit.default, exp_search_fd.default, exp_authorization_tg.default, exp_ds_cat_diversity.test, exp_main_page_inapp_message_vs_fullsreen.test",
		"traceparent":         "00-e13ddedb8cef1509ffdde320cdfbf217-83ff177c3dd7d64f-01",
		"x-trace-id":          "e13ddedb8cef1509ffdde320cdfbf217",
		"x-span-id":           "83ff177c3dd7d64f",
		"Cookie":              cookie,
		"Connection":          "keep-alive",
	}
}
