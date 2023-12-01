import { handleErrors } from '../../../assets/helper/handleErrors';

export const load = async (context: { params: { id: string } }): Promise<object> => {
    const {id} = context.params;
    const apiUrl: string = `/api/v1/user/${id}`;

    return fetch(apiUrl)
        .then(handleErrors)
        .then(merchant => {
            return {
                merchant,
                metaTitle: 'Supermarkt: ' + merchant?.name,
            };
        })
        .catch(error => {
            console.error("Failed to fetch shopping lists data:", error.message);
            return {
                merchant: [],
                metaTitle: 'Error',
            };
        });
};

