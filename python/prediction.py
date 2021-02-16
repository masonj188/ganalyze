import ember
import lightgbm as lgb
import sys

def main():
    #lgbm_model = lgb.Booster(model_file="ember_model_2017.txt")
    lgbm_model = lgb.Booster(model_file="ember_model_2018.txt")

    try:
        file = sys.argv[1]
        if file.endswith('.exe') or file.endswith('.dll'):
            test_data = open(file, "rb").read()
        else:
            print(-1, file=sys.stderr)
            return
        
    except IndexError:
        print(-1, file=sys.stderr)
        return
    
    print(int(round(ember.predict_sample(lgbm_model, test_data))), file=sys.stderr)

if __name__ == "__main__":
    main()
