token = request.query_params.get("token") or request.data.get("token")
        data = {'request': 'JumpValidateRequest',
                'platformCode': platformCode,
                'token': urllib.parse.quote(token)}

        resp = requests.post(urljoin(pms, "/api/user/validate"), json=data)
        result = resp.json()